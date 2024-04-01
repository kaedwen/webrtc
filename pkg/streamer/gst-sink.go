package streamer

import (
	"fmt"

	"github.com/go-gst/go-gst/gst"
	"github.com/go-gst/go-gst/gst/app"
	"github.com/kaedwen/webrtc/pkg/common"
	"github.com/pion/webrtc/v3/pkg/media"
	"go.uber.org/zap"
)

func setCallback(sink *app.Sink, ch chan media.Sample) {
	sink.SetCallbacks(&app.SinkCallbacks{
		NewSampleFunc: func(sink *app.Sink) gst.FlowReturn {
			// Pull the sample that triggered this callback
			sample := sink.PullSample()
			if sample == nil {
				return gst.FlowEOS
			}

			// Retrieve the buffer from the sample
			buffer := sample.GetBuffer()
			if buffer == nil {
				return gst.FlowError
			}

			// At this point, buffer is only a reference to an existing memory region somewhere.
			// When we want to access its content, we have to map it while requesting the required
			// mode of access (read, read/write).
			data := buffer.Map(gst.MapRead).AsUint8Slice()
			defer buffer.Unmap()

			ch <- media.Sample{Data: data, Duration: *buffer.Duration().AsDuration()}

			return gst.FlowOK
		},
	})
}

func CreateVideoPipelineSinkWithLaunch(lg *zap.Logger, s StreamElement) (*gst.Pipeline, <-chan media.Sample, error) {
	pb := NewPipelineBuilder()

	pb.AddWithProperties("v4l2src", s.Properties)
	pb.AddFilter(s.SrcCaps)
	pb.Add("videoconvert")

	if s.Queue {
		pb.Add("queue")
	}

	switch s.Codec {
	case common.VP8:
		pb.AddWithProperties("vp8enc", map[string]any{
			"bitrate":           s.Bitrate,
			"error-resilient":   "partitions",
			"keyframe-max-dist": int(10),
			"cpu-used":          int(5),
			"deadline":          int(1),
			"auto-alt-ref":      true,
		})
	case common.VP9:
		pb.AddWithProperties("vp9enc", map[string]any{
			"bitrate": s.Bitrate,
		})
	case common.H264:
		pb.AddFilter(NewCaps("video/x-raw", map[string]any{
			"format": "I420",
		}))
		pb.AddWithProperties("x264enc", map[string]any{
			"bitrate":      s.Bitrate,
			"speed-preset": "ultrafast",
			"tune":         "zerolatency",
			"key-int-max":  int(20),
		})
		pb.AddFilter(NewCaps("video/x-h264", map[string]any{
			"stream-format": "byte-stream",
		}))
	default:
		return nil, nil, fmt.Errorf("unsupported video codec given - %s", s.Codec)
	}

	pb.AddWithProperties("appsink", map[string]any{
		"name": "appsink",
	})

	ps := pb.Build()
	lg.Info("launch pipeline", zap.String("definition", ps))

	pipeline, err := gst.NewPipelineFromString(ps)
	if err != nil {
		return nil, nil, err
	}

	elem, err := pipeline.GetElementByName("appsink")
	if err != nil {
		return nil, nil, err
	}

	appsink := app.SinkFromElement(elem)
	ch := make(chan media.Sample, 100)
	setCallback(appsink, ch)

	return pipeline, ch, nil
}

func CreateVideoPipelineSink(lg *zap.Logger, s StreamElement) (*gst.Pipeline, <-chan media.Sample, error) {
	// Create a pipeline
	pipeline, err := gst.NewPipeline("pion-video-pipeline")
	if err != nil {
		return nil, nil, err
	}

	// Create the src
	src, err := gst.NewElement("v4l2src")
	if err != nil {
		return nil, nil, err
	}

	src.Set("device", "/dev/video4")

	src_c := gst.NewEmptySimpleCaps("video/x-raw")
	src_c.SetValue("width", int(320))
	src_c.SetValue("height", int(240))
	lg.Info("capsfilter", zap.String("caps", src_c.String()))

	// just to be on the save side
	conv, err := gst.NewElement("videoconvert")
	if err != nil {
		return nil, nil, err
	}

	// add a queue
	queue, err := gst.NewElement("queue")
	if err != nil {
		return nil, nil, err
	}

	enc_in_c := gst.NewEmptySimpleCaps("video/x-raw")
	enc_in_c.SetValue("format", "I420")
	lg.Info("capsfilter", zap.String("caps", enc_in_c.String()))

	// Create the enc
	enc, err := gst.NewElement("x264enc")
	if err != nil {
		return nil, nil, err
	}

	enc.SetProperty("speed-preset", "ultrafast")
	enc.SetProperty("tune", "zerolatency")
	enc.SetProperty("key-int-max", 20)
	//enc.SetProperty("bitrate", 300)

	enc_out_c := gst.NewEmptySimpleCaps("video/x-h264")
	enc_out_c.SetValue("stream-format", "byte-stream")
	lg.Info("capsfilter", zap.String("caps", enc_out_c.String()))

	// Create the sink
	appsink, err := app.NewAppSink()
	if err != nil {
		return nil, nil, err
	}

	ch := make(chan media.Sample, 100)
	setCallback(appsink, ch)

	// Add the elements to the pipeline
	err = pipeline.AddMany(src, conv, queue, enc, appsink.Element)
	if err != nil {
		return nil, nil, err
	}

	// link the elements
	err = src.LinkFiltered(conv, src_c)
	if err != nil {
		return nil, nil, err
	}

	err = conv.Link(queue)
	if err != nil {
		return nil, nil, err
	}

	err = queue.LinkFiltered(enc, enc_in_c)
	if err != nil {
		return nil, nil, err
	}

	err = enc.LinkFiltered(appsink.Element, enc_out_c)
	if err != nil {
		return nil, nil, err
	}

	return pipeline, ch, nil
}

func CreateAudioPipelineSink(s StreamElement, lg *zap.Logger) (*gst.Pipeline, <-chan media.Sample, error) {
	// Create a pipeline
	pipeline, err := gst.NewPipeline("pion-audio-pipeline")
	if err != nil {
		return nil, nil, err
	}

	elems := make(ElementList, 0)

	// Create the src
	src, err := gst.NewElement(s.Kind)
	if err != nil {
		return nil, nil, err
	}

	for name, value := range s.Properties {
		src.Set(name, value)
	}

	// if s.Caps != nil {
	// 	c := s.Caps.Build()
	// 	lg.Info("capsfilter", zap.String("caps", c.String()))
	// 	elems = append(elems, Element{src, c})
	// } else {
	// 	elems = append(elems, Element{src, nil})
	// }

	// just to be on the save side
	conv, err := gst.NewElement("audioconvert")
	if err != nil {
		return nil, nil, err
	}
	elems = append(elems, Element{conv, nil})

	if s.Queue {
		// add a queue
		queue, err := gst.NewElement("queue")
		if err != nil {
			return nil, nil, err
		}
		elems = append(elems, Element{queue, nil})
	}

	switch s.Codec {
	case common.OPUS:
		// Create the enc
		enc, err := gst.NewElement("opusenc")
		if err != nil {
			return nil, nil, err
		}
		elems = append(elems, Element{enc, nil})
	default:
		return nil, nil, fmt.Errorf("unsupported audio codec given - %s", s.Codec)
	}

	// Create the sink
	appsink, err := app.NewAppSink()
	if err != nil {
		return nil, nil, err
	}
	elems = append(elems, Element{appsink.Element, nil})

	ch := make(chan media.Sample, 100)
	setCallback(appsink, ch)

	// Add the elements to the pipeline
	err = pipeline.AddMany(elems.List()...)
	if err != nil {
		return nil, nil, err
	}

	// link the elements
	err = elems.Link()
	if err != nil {
		return nil, nil, err
	}

	return pipeline, ch, nil
}

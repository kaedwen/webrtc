package streamer

import (
	"github.com/pion/webrtc/v3/pkg/media"
	"github.com/tinyzimmer/go-gst/gst"
	"github.com/tinyzimmer/go-gst/gst/app"
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

			ch <- media.Sample{Data: data, Duration: buffer.Duration()}

			return gst.FlowOK
		},
	})
}

func CreateVideoPipelineSink(s StreamElement) (*gst.Pipeline, <-chan media.Sample, error) {
	// Create a pipeline
	pipeline, err := gst.NewPipeline("pion-video-pipeline")
	if err != nil {
		return nil, nil, err
	}

	elems := make([]*gst.Element, 0)

	// Create the src
	src, err := gst.NewElement(s.Kind)
	if err != nil {
		return nil, nil, err
	}
	elems = append(elems, src)

	for name, value := range s.Properties {
		src.Set(name, value)
	}

	if s.Caps != nil {
		filter, err := gst.NewElement("capsfilter")
		if err != nil {
			return nil, nil, err
		}

		filter.Set("caps", s.ToGstCaps())
		elems = append(elems, filter)
	}

	// just to be on the save side
	conv, err := gst.NewElement("videoconvert")
	if err != nil {
		return nil, nil, err
	}
	elems = append(elems, conv)

	if s.Queue {
		// add a queue
		queue, err := gst.NewElement("queue")
		if err != nil {
			return nil, nil, err
		}
		elems = append(elems, queue)
	}

	// Create the enc
	enc, err := gst.NewElement("vp8enc")
	if err != nil {
		return nil, nil, err
	}
	elems = append(elems, enc)

	enc.SetProperty("error-resilient", "partitions")
	enc.SetProperty("keyframe-max-dist", 10)
	enc.SetProperty("auto-alt-ref", true)
	enc.SetProperty("cpu-used", 5)
	enc.SetProperty("deadline", 1)

	// Create the sink
	appsink, err := app.NewAppSink()
	if err != nil {
		return nil, nil, err
	}
	elems = append(elems, appsink.Element)

	ch := make(chan media.Sample, 100)
	setCallback(appsink, ch)

	// Add the elements to the pipeline
	pipeline.AddMany(elems...)

	// link the elements
	gst.ElementLinkMany(elems...)

	return pipeline, ch, nil
}

func CreateAudioPipelineSink(s StreamElement) (*gst.Pipeline, <-chan media.Sample, error) {
	// Create a pipeline
	pipeline, err := gst.NewPipeline("pion-audio-pipeline")
	if err != nil {
		return nil, nil, err
	}

	elems := make([]*gst.Element, 0)

	// Create the src
	src, err := gst.NewElement(s.Kind)
	if err != nil {
		return nil, nil, err
	}
	elems = append(elems, src)

	for name, value := range s.Properties {
		src.Set(name, value)
	}

	if s.Caps != nil {
		filter, err := gst.NewElement("capsfilter")
		if err != nil {
			return nil, nil, err
		}

		filter.Set("caps", s.ToGstCaps())
		elems = append(elems, filter)
	}

	// just to be on the save side
	conv, err := gst.NewElement("audioconvert")
	if err != nil {
		return nil, nil, err
	}
	elems = append(elems, conv)

	if s.Queue {
		// add a queue
		queue, err := gst.NewElement("queue")
		if err != nil {
			return nil, nil, err
		}
		elems = append(elems, queue)
	}

	// Create the enc
	enc, err := gst.NewElement("opusenc")
	if err != nil {
		return nil, nil, err
	}
	elems = append(elems, enc)

	// Create the sink
	appsink, err := app.NewAppSink()
	if err != nil {
		return nil, nil, err
	}
	elems = append(elems, appsink.Element)

	ch := make(chan media.Sample, 100)
	setCallback(appsink, ch)

	// Add the elements to the pipeline
	pipeline.AddMany(elems...)

	// link the elements
	gst.ElementLinkMany(elems...)

	return pipeline, ch, nil
}

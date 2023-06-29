package webrtc

import (
	"github.com/pion/webrtc/v3/pkg/media"
	"github.com/tinyzimmer/go-gst/gst"
	"github.com/tinyzimmer/go-gst/gst/app"
)

func init() {
	gst.Init(nil)
}

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

func CreateVideoPipeline(s string) (*gst.Pipeline, <-chan media.Sample, error) {
	// Create a pipeline
	pipeline, err := gst.NewPipeline("")
	if err != nil {
		return nil, nil, err
	}

	// Create the src
	src, err := gst.NewElement(s)
	if err != nil {
		return nil, nil, err
	}

	// Create the enc
	enc, err := gst.NewElement("vp8enc")
	if err != nil {
		return nil, nil, err
	}

	enc.SetProperty("error-resilient", "partitions")
	enc.SetProperty("keyframe-max-dist", 10)
	enc.SetProperty("auto-alt-ref", true)
	enc.SetProperty("cpu-used", 5)
	enc.SetProperty("deadline", 1)

	// Create the sink
	sink, err := app.NewAppSink()
	if err != nil {
		return nil, nil, err
	}

	ch := make(chan media.Sample, 100)
	setCallback(sink, ch)

	// Add the elements to the pipeline
	pipeline.AddMany(src, enc, sink.Element)

	// link the elements
	gst.ElementLinkMany(src, enc, sink.Element)

	return pipeline, ch, nil
}

func CreateAudioPipeline(s string) (*gst.Pipeline, <-chan media.Sample, error) {
	// Create a pipeline
	pipeline, err := gst.NewPipeline("")
	if err != nil {
		return nil, nil, err
	}

	// Create the src
	src, err := gst.NewElement(s)
	if err != nil {
		return nil, nil, err
	}

	// Create the enc
	enc, err := gst.NewElement("opusenc")
	if err != nil {
		return nil, nil, err
	}

	// Create the sink
	sink, err := app.NewAppSink()
	if err != nil {
		return nil, nil, err
	}

	ch := make(chan media.Sample, 100)
	setCallback(sink, ch)

	// Add the elements to the pipeline
	pipeline.AddMany(src, enc, sink.Element)

	// link the elements
	gst.ElementLinkMany(src, enc, sink.Element)

	return pipeline, ch, nil
}

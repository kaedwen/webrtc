package streamer

import (
	"fmt"

	"github.com/tinyzimmer/go-gst/gst"
	"github.com/tinyzimmer/go-gst/gst/app"
)

type SrcPipeline struct {
	*gst.Pipeline
	src *app.Source
}

func CreateAudioPipelineSrc(dst StreamElement) (*SrcPipeline, error) {
	// Create a pipeline
	pipeline, err := gst.NewPipeline("")
	if err != nil {
		return nil, err
	}

	elems := make([]*gst.Element, 0)

	// Create the src
	appsrc, err := app.NewAppSrc()
	if err != nil {
		return nil, err
	}
	elems = append(elems, appsrc.Element)

	appsrc.SetFormat(gst.FormatTime)
	appsrc.SetDoTimestamp(true)
	appsrc.SetLive(true)

	// Create the opus decoder
	decCodec, err := gst.NewElement("rtpopusdepay")
	if err != nil {
		return nil, err
	}
	elems = append(elems, decCodec)

	// Create the bin decoder
	decBin, err := gst.NewElement("decodebin")
	if err != nil {
		return nil, err
	}
	elems = append(elems, decBin)

	// Create the sink
	sink, err := gst.NewElement(dst.Kind)
	if err != nil {
		return nil, err
	}
	elems = append(elems, sink)

	for name, value := range dst.Properties {
		sink.SetProperty(name, value)
	}

	// Add the elements to the pipeline
	pipeline.AddMany(elems...)

	// link the elements
	gst.ElementLinkMany(elems...)

	return &SrcPipeline{
		Pipeline: pipeline,
		src:      appsrc,
	}, nil
}

func (p *SrcPipeline) Push(data []byte) error {
	err := p.src.PushBuffer(gst.NewBufferFromBytes(data))
	if err != gst.FlowOK {
		return fmt.Errorf("failed to send bytes to gst pipeline - %s", err.String())
	}

	return nil
}

package streamer

import (
	"fmt"

	"github.com/go-gst/go-gst/gst"
	"github.com/go-gst/go-gst/gst/app"
)

type SrcPipeline struct {
	*gst.Pipeline
	src *app.Source
}

func NewSrcPipeline(p *gst.Pipeline, src *app.Source) *SrcPipeline {
	return &SrcPipeline{
		Pipeline: p,
		src:      src,
	}
}

func CreateAudioPipelineSrc(dst StreamElement) (*SrcPipeline, error) {
	// Create a pipeline
	pipeline, err := gst.NewPipeline("")
	if err != nil {
		return nil, err
	}

	elems, err := gst.NewElementMany("appsrc", "rtpjitterbuffer", "rtpopusdepay", "opusdec", dst.Kind)
	if err != nil {
		return nil, err
	}

	caps := gst.NewEmptySimpleCaps("application/x-rtp")
	caps.SetValue("media", "audio")
	caps.SetValue("clock-rate", 48000)
	caps.SetValue("payload", 96)
	caps.SetValue("encoding-name", "OPUS")

	appsrc := app.SrcFromElement(elems[0])
	appsrc.SetFormat(gst.FormatTime)
	appsrc.SetDoTimestamp(true)
	appsrc.SetLive(true)
	appsrc.SetCaps(caps)

	// Create the sink
	sink := elems[len(elems)-1]

	for name, value := range dst.Properties {
		sink.Set(name, value)
	}

	// Add the elements to the pipeline and link them
	err = pipeline.AddMany(elems...)
	if err != nil {
		return nil, err
	}
	err = gst.ElementLinkMany(elems...)
	if err != nil {
		return nil, err
	}

	return NewSrcPipeline(pipeline, appsrc), nil
}

func (p *SrcPipeline) Push(data []byte) error {
	err := p.src.PushBuffer(gst.NewBufferFromBytes(data))
	if err != gst.FlowOK {
		return fmt.Errorf("failed to send bytes to gst pipeline - %s", err.String())
	}

	return nil
}

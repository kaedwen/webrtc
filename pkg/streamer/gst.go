package streamer

import (
	"github.com/go-gst/go-gst/gst"
	"github.com/go-gst/go-gst/gst/app"
	"go.uber.org/zap"
)

func init() {
	gst.Init(nil)
}

type Element struct {
	*gst.Element
	filter *gst.Caps
}

type ElementList []Element

func (elems ElementList) Link() error {
	for idx, elem := range elems {
		if idx == 0 {
			// skip the first one as the loop always links previous to current
			continue
		}
		pe := elems[idx-1]
		if elem.filter != nil {
			if err := pe.LinkFiltered(elem.Element, pe.filter); err != nil {
				return err
			}
		} else {
			if err := pe.Link(elem.Element); err != nil {
				return err
			}
		}

	}
	return nil
}

func (elems ElementList) List() []*gst.Element {
	l := make([]*gst.Element, 0, len(elems))
	for _, e := range elems {
		l = append(l, e.Element)
	}
	return l
}

func handleMessage(msg *gst.Message) error {
	switch msg.Type() {
	case gst.MessageEOS:
		return app.ErrEOS
	case gst.MessageError:
		return msg.ParseError()
	}

	return nil
}

func LoopBus(lg *zap.Logger, pipeline *gst.Pipeline) {
	// Retrieve the bus from the pipeline
	bus := pipeline.GetPipelineBus()

	// Loop over messsages from the pipeline
	go func() {
		for {
			msg := bus.TimedPop(gst.ClockTimeNone)
			if msg == nil {
				return
			}
			if err := handleMessage(msg); err != nil {
				lg.Error("failed to handle message", zap.Error(err))
			}
		}
	}()
}

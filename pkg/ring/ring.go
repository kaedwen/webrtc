package ring

import (
	"context"
	"fmt"
	"net/url"

	"github.com/holoplot/go-evdev"
	"github.com/kaedwen/webrtc/pkg/common"
	"github.com/kaedwen/webrtc/pkg/ring/sonos"
	"go.uber.org/zap"
)

type PlayHandler interface {
	Play(*url.URL) error
	Watch(context.Context) error
}

type RingHandler struct {
	lg           *zap.Logger
	cfg          *common.ConfigRing
	playHandlers []PlayHandler
}

func NewRingHandler(lg *zap.Logger, cfg *common.ConfigRing) (*RingHandler, error) {
	spl, err := sonos.NewSonosHandler(lg.With(zap.String("context", "sonos")), cfg)
	if err != nil {
		return nil, err
	}

	return &RingHandler{lg, cfg, []PlayHandler{spl}}, nil
}

func (h *RingHandler) Watch(ctx context.Context) error {
	if h.cfg.Device == nil {
		h.lg.Warn("nothing to watch for key press")
		return nil
	}

	if h.cfg.JingleBaseUri == nil {
		h.lg.Warn("missing jingle base uri, nothing to play")
		return nil
	}

	for _, p := range h.playHandlers {
		if err := p.Watch(ctx); err != nil {
			return err
		}
	}

	d, err := evdev.Open(*h.cfg.Device)
	if err != nil {
		return err
	}

	vMajor, vMinor, vMicro := d.DriverVersion()
	h.lg.Info("input driver running", zap.String("version", fmt.Sprintf("%d.%d.%d", vMajor, vMinor, vMicro)))

	key := evdev.KEYFromString[h.cfg.Key]

	bu, err := url.Parse(*h.cfg.JingleBaseUri)
	if err != nil {
		return err
	}

	go func() {
		err := d.NonBlock()
		if err != nil {
			h.lg.Fatal("failed to set non block", zap.Error(err))
		}

		for {
			e, err := d.ReadOne()
			if err == nil && e.Code == key && e.Value == 0 {
				tu := bu.JoinPath(h.cfg.JinglePath)
				for _, p := range h.playHandlers {
					if err := p.Play(tu); err != nil {
						h.lg.Error("failed to play", zap.Error(err))
					}
				}
			}
		}
	}()

	go func() {
		<-ctx.Done()
		_ = d.Close()
	}()

	return nil
}

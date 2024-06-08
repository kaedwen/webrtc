package ring

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/holoplot/go-evdev"
	"github.com/kaedwen/webrtc/pkg/common"
	"github.com/kaedwen/webrtc/pkg/ring/sonos"
	"go.uber.org/zap"
)

type PlayHandler interface {
	Play(context.Context, *url.URL) error
	Watch(context.Context) error
}

type RingHandler struct {
	lg           *zap.Logger
	cfg          *common.ConfigRing
	playHandlers []PlayHandler
}

func NewRingHandler(ctx context.Context, lg *zap.Logger, cfg *common.ConfigRing) error {
	spl, err := sonos.NewSonosHandler(lg.With(zap.String("context", "sonos")), cfg)
	if err != nil {
		return err
	}

	rh := &RingHandler{lg, cfg, []PlayHandler{spl}}

	if err = rh.watch(ctx); err != nil {
		return err
	}

	return nil
}

func (h *RingHandler) watch(ctx context.Context) error {
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
				// run all players
				tu := bu.JoinPath(h.cfg.JinglePath.String())
				for _, p := range h.playHandlers {
					if err := p.Play(ctx, tu); err != nil {
						h.lg.Error("failed to play", zap.Error(err))
					}
				}

				// run webhooks when configured
				if h.cfg.HomeassistantWebhook != nil {
					h.lg.Info("Triggering Homeassistant Webhook", zap.String("hook", *h.cfg.HomeassistantWebhook))
					ctxt, cancel := context.WithTimeout(ctx, 30*time.Second)

					err := h.TriggerWebhook(ctxt, *h.cfg.HomeassistantWebhook)
					if err != nil {
						h.lg.Error("failed to trigger webhook", zap.Error(err))
					}

					cancel()
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

func (h *RingHandler) TriggerWebhook(ctx context.Context, hook string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, *h.cfg.HomeassistantWebhook, nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("received wrong status code - %d", res.StatusCode)
	}

	return nil
}

package sonos

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/hashicorp/mdns"
	"github.com/kaedwen/webrtc/pkg/common"
	"go.uber.org/zap"
)

type sonosInfo struct {
	Device struct {
		Id               string   `json:"id"`
		SerialNumber     string   `json:"serialNumber"`
		Model            string   `json:"model"`
		ModelDisplayName string   `json:"modelDisplayName"`
		Color            string   `json:"color"`
		Capabilities     []string `json:"capabilities"`
		ApiVersion       string   `json:"apiVersion"`
		MinApiVersion    string   `json:"minApiVersion"`
		Name             string   `json:"name"`
		WebsocketUrl     string   `json:"websocketUrl"`
		SoftwareVersion  string   `json:"softwareVersion"`
		HwVersion        string   `json:"hwVersion"`
		SwGen            int      `json:"swGen"`
	} `json:"device"`
	HouseholdId  string `json:"householdId"`
	PlayerId     string `json:"playerId"`
	GroupId      string `json:"groupId"`
	WebsocketUrl string `json:"websocketUrl"`
	RestUrl      string `json:"restUrl"`
}

type sonosAudioClip struct {
	Name      string `json:"name"`
	AppId     string `json:"appId"`
	ClipType  string `json:"clipType,omitempty"`
	StreamUrl string `json:"streamUrl"`
	Volume    int    `json:"volume"`
}

type SonosHandler struct {
	lg      *zap.Logger
	cfg     *common.ConfigRing
	players map[string]*SonosPlayer
}

type SonosPlayer struct {
	client  *http.Client
	address *url.URL
	info    sonosInfo
}

func NewSonosPlayer(address *url.URL) (*SonosPlayer, error) {
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}

	return &SonosPlayer{client, address, sonosInfo{}}, nil
}

func NewSonosHandler(lg *zap.Logger, cfg *common.ConfigRing) (*SonosHandler, error) {
	return &SonosHandler{lg, cfg, make(map[string]*SonosPlayer)}, nil
}

func (p *SonosPlayer) init(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.address.JoinPath("api/v1/players/local/info").String(), nil)
	if err != nil {
		return err
	}

	// for some reason there must be a api key given
	// does not matter what
	req.Header.Add("X-Sonos-Api-Key", uuid.NewString())

	res, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("received status %d", res.StatusCode)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &p.info)
	if err != nil {
		return err
	}

	return nil
}

func (p *SonosPlayer) Play(ctx context.Context, uri *url.URL, volume int) error {
	sab := sonosAudioClip{
		Name:      "Pull Bell",
		AppId:     "com.acme.app",
		ClipType:  "CHIME",
		StreamUrl: uri.String(),
		Volume:    volume,
	}

	b, err := json.Marshal(sab)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.address.JoinPath(fmt.Sprintf("api/v1/players/%s/audioClip", p.info.PlayerId)).String(), bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	// for some reason there must be a api key given
	// does not matter what
	req.Header.Add("X-Sonos-Api-Key", uuid.NewString())
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	res, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("received status %d", res.StatusCode)
	}

	return nil
}

func (h *SonosHandler) Watch(ctx context.Context) error {
	// Make a channel for results and start listening
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case e := <-entriesCh:
				h.lg.Info("found player", zap.String("name", e.Name), zap.String("address", e.AddrV4.String()))
				if _, has := h.players[e.Name]; !has {
					p, err := NewSonosPlayer(&url.URL{
						Scheme: "https",
						Host:   net.JoinHostPort(e.AddrV4.String(), "1443"),
					})
					if err != nil {
						h.lg.Error("failed to create player", zap.Error(err))
						continue
					}
					if err := p.init(ctx); err != nil {
						h.lg.Error("failed to init player", zap.Error(err))
						continue
					}

					if h.cfg.SonosTarget != "-" && p.info.Device.Name != h.cfg.SonosTarget {
						h.lg.Info("skip player", zap.String("name", p.info.Device.Name))
						continue
					}

					h.players[e.Name] = p
				}
			}
		}
	}()

	params := mdns.DefaultParams("_sonos._tcp")
	params.Entries = entriesCh

	if h.cfg.NoIPv6 {
		params.DisableIPv6 = true
	}

	// Start the lookup
	return mdns.Query(params)
}

func (h *SonosHandler) Play(ctx context.Context, uri *url.URL) error {
	for _, p := range h.players {
		h.lg.Info("playing clip", zap.String("target", p.address.String()), zap.String("clip", uri.String()))
		if err := p.Play(ctx, uri, h.cfg.SonosVolume); err != nil {
			h.lg.Error("failed to play", zap.String("address", p.address.String()), zap.Error(err))
		}
	}

	return nil
}

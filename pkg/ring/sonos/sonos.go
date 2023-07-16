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
	StreamUrl string `json:"streamUrl"`
	Volume    int    `json:"volume"`
}

type SonosHandler struct {
	lg      *zap.Logger
	target  string
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

func NewSonosHandler(lg *zap.Logger, target string) (*SonosHandler, error) {
	return &SonosHandler{lg, target, make(map[string]*SonosPlayer)}, nil
}

func (p *SonosPlayer) init() error {
	req, err := http.NewRequest(http.MethodGet, p.address.JoinPath("api/v1/players/local/info").String(), nil)
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

func (p *SonosPlayer) Play(uri *url.URL) error {
	sab := sonosAudioClip{
		Name:      "Pull Bell",
		AppId:     "com.acme.app",
		StreamUrl: uri.String(),
		Volume:    5,
	}

	b, err := json.Marshal(sab)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, p.address.JoinPath(fmt.Sprintf("api/v1/players/%s/audioClip", p.info.PlayerId)).String(), bytes.NewBuffer(b))
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
					if err := p.init(); err != nil {
						h.lg.Error("failed to init player", zap.Error(err))
						continue
					}

					h.players[e.Name] = p
				}
			}
		}
	}()

	// Start the lookup
	return mdns.Lookup("_sonos._tcp", entriesCh)
}

func (h *SonosHandler) Play(uri *url.URL) error {
	for _, p := range h.players {
		if err := p.Play(uri); err != nil {
			h.lg.Error("failed to play", zap.String("address", p.address.String()), zap.Error(err))
		}
	}

	return nil
}

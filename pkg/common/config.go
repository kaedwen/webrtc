package common

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexflint/go-arg"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type File = ConfigFile
type VideoSrc = ConfigVideoSourceStream
type AudioSrc = ConfigAudioSourceStream
type AudioSink = ConfigAudioSinkStream
type Logging = ConfigLogging
type Ring = ConfigRing
type Http = ConfigHTTP

type Config struct {
	File
	VideoSrc  `yaml:"video-src"`
	AudioSrc  `yaml:"audio-src"`
	AudioSink `yaml:"audio-sink"`
	Logging   `yaml:"logging"`
	Ring      `yaml:"ring"`
	Http      `yaml:"http"`
}

type Path struct {
	string
}

func (p Path) String() string {
	return p.string
}

type ConfigStream struct {
	VideoSrc  ConfigVideoSourceStream // video src for webrtc send
	AudioSrc  ConfigAudioSourceStream // audio src for webrtc send
	AudioSink ConfigAudioSinkStream   // audio sink for webrtc receive
}

type ConfigFile struct {
	Path *Path `arg:"--config"`
}

type ConfigLogging struct {
	Level       zap.AtomicLevel `arg:"--log-level,env:LOG_LEVEL" yaml:"level" default:"debug"`
	Development bool            `arg:"--log-development,env:LOG_DEVELOPMENT" yaml:"development"`
}

type ConfigRing struct {
	Device               *string `arg:"--input-device,env:INPUT_DEVICE" yaml:"input"`
	Key                  string  `arg:"--ring-key" default:"KEY_F1" yaml:"key"`
	JingleBaseUri        *string `arg:"--jingle-base-uri,env:JINGLE_BASE_URI" yaml:"jingle-base-uri"`
	JinglePath           Path    `arg:"--jingle-path,env:JINGLE_PATH" default:"audio/ding-dong.wav" yaml:"jingle-path"`
	SonosTarget          string  `arg:"--sonos-target,env:SONOS_TARGET" yaml:"sonos-target" default:"-"`
	SonosVolume          int     `arg:"--sonos-volume,env:SONOS_VOLUME" yaml:"sonos-volume" default:"50"`
	HomeassistantWebhook *string `arg:"--ha-webhook,env:HA_WEBHOOK" yaml:"ha-webhook"`
	NoIPv6               bool    `arg:"--disable-ipv6,env:NO_IPV6" yaml:"no-ipv6"`
}

type ConfigHTTP struct {
	Host             string  `arg:"--http-host,env:HTTP_HOST" yaml:"host" default:"0.0.0.0"`
	Port             uint    `arg:"--http-port,env:HTTP_PORT" yaml:"port" default:"8080"`
	Tls              bool    `arg:"--http-tls,env:HTTP_TLS" yaml:"tls" default:"true"`
	TlsKey           *string `arg:"--http-tls-key,env:HTTP_TLS_KEY" yaml:"tls-key"`
	TlsCert          *string `arg:"--http-tls-cert,env:HTTP_TLS_CERT" yaml:"tls-cert"`
	PathGetLiveness  string  `arg:"env:HTTP_PATH_LIVENESS" yaml:"liveness" default:"/healthz"`
	PathGetReadiness string  `arg:"env:HTTP_PATH_READINESS" yaml:"readiness" default:"/readyz"`
	StaticPath       *Path   `arg:"--http-static,env:HTTP_STATIC" yaml:"static"`
}

type ConfigVideoSourceStream struct {
	Source    string      `arg:"--video-src,env:VIDEO_SRC" yaml:"source" default:"v4l2src"`
	Device    string      `arg:"--video-src-device,env:VIDEO_SRC_DEVICE" yaml:"device" default:"/dev/video0"`
	Codec     StreamCodec `arg:"--video-src-codec,env:VIDEO_SRC_CODEC" yaml:"codec" default:"vp8"`
	Height    uint        `arg:"--video-src-height,env:VIDEO_SRC_HEIGHT" yaml:"height" default:"480"`
	Width     uint        `arg:"--video-src-width,env:VIDEO_SRC_WIDTH" yaml:"width" default:"640"`
	Framerate uint        `arg:"--video-src-fps,env:VIDEO_SRC_FPS" yaml:"fps" default:"30"`
	Bitrate   uint        `arg:"--video-src-bps,env:VIDEO_SRC_BPS" yaml:"bps" default:"300"`
	Queue     bool        `arg:"--video-src-queue,env:VIDEO_SRC_QUEUE" yaml:"queue" default:"false"`
}

type ConfigAudioSourceStream struct {
	Source     string      `arg:"--audio-src,env:AUDIO_SRC" yaml:"source" default:"alsasrc"`
	DeviceName string      `arg:"--audio-src-device-name,env:AUDIO_SRC_DEVICE_NAME" yaml:"name" default:"default"`
	Device     *string     `arg:"--audio-src-device,env:AUDIO_SRC_DEVICE" yaml:"device"`
	Codec      StreamCodec `arg:"--audio-src-codec,env:AUDIO_SRC_CODEC" yaml:"codec" default:"opus"`
	Channels   uint        `arg:"--audio-src-channels,env:AUDIO_SRC_CHANNELS" yaml:"channels" default:"1"`
	Queue      bool        `arg:"--audio-src-queue,env:AUDIO_SRC_QUEUE" yaml:"queue" default:"false"`
}

type ConfigAudioSinkStream struct {
	Sink       string  `arg:"--audio-sink,env:AUDIO_SINK" yaml:"sink" default:"alsasink"`
	DeviceName string  `arg:"--audio-sink-device-name,env:AUDIO_SINK_DEVICE_NAME" yaml:"name" default:"default"`
	Device     *string `arg:"--audio-sink-device,env:AUDIO_SINK_DEVICE" yaml:"device"`
	Codec      string  `arg:"--audio-sink-codec,env:AUDIO_SINK_CODEC" yaml:"codec" default:"opus"`
	Channels   uint    `arg:"--audio-sink-channels,env:AUDIO_SINK_CHANNELS" yaml:"channels" default:"1"`
	Queue      bool    `arg:"--audio-sink-queue,env:AUDIO_SINK_QUEUE" yaml:"queue" default:"false"`
}

func (c *Config) Stream() *ConfigStream {
	return &ConfigStream{
		VideoSrc:  c.VideoSrc,
		AudioSrc:  c.AudioSrc,
		AudioSink: c.AudioSink,
	}
}

func (c *ConfigHTTP) Address() string {
	return net.JoinHostPort(c.Host, fmt.Sprint(c.Port))
}

func (p *Path) UnmarshalText(b []byte) error {
	p.string = string(b)

	if sc, had := strings.CutPrefix(p.string, "~/"); had {
		dir, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		p.string = filepath.Join(dir, sc)
	}

	p.string = os.ExpandEnv(p.string)

	return nil
}

func (c *Config) MustParse() {
	// first parse just the config file flag
	arg.MustParse(&c.File)

	if c.File.Path != nil {
		if _, err := os.Stat(c.File.Path.String()); err == nil {
			if data, err := os.ReadFile(c.File.Path.String()); err == nil {
				err := yaml.Unmarshal(data, c)
				if err != nil {
					panic(err)
				}
			}
		}
	}

	// parse the rest
	arg.MustParse(c)
}

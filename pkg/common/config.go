package common

import (
	"fmt"
	"net"

	"github.com/alexflint/go-arg"
)

type Config struct {
	VideoOut ConfigVideoOutputStream `arg:"group:VideoOut"`
	AudioOut ConfigAudioOutputStream `arg:"group:AudioOut"`
	AudioIn  ConfigAudioInputStream  `arg:"group:AudioIn"`
	Logging  ConfigLogging           `arg:"group:Logging"`
	Ring     ConfigRing              `arg:"group:Ring"`
	Http     ConfigHTTP              `arg:"group:Http"`
}

type ConfigStream struct {
	VideoOut ConfigVideoOutputStream
	AudioOut ConfigAudioOutputStream
	AudioIn  ConfigAudioInputStream
}

type ConfigLogging struct {
	Level string `arg:"--log-level,env:LOG_LEVEL" default:"debug"`
}

type ConfigRing struct {
	Device               *string `arg:"--input-device,env:INPUT_DEVICE"`
	Key                  string  `arg:"--ring-key" default:"KEY_F1"`
	JingleBaseUri        *string `arg:"--jingle-base-uri,env:JINGLE_BASE_URI"`
	JinglePath           string  `arg:"--jingle-path,env:JINGLE_PATH" default:"audio/ding-dong.wav"`
	SonosTarget          string  `arg:"--sonos-target,env:SONOS_TARGET" default:"-"`
	SonosVolume          int     `arg:"--sonos-volume,env:SONOS_VOLUME" default:"50"`
	HomeassistantWebhook *string `arg:"--ha-webhook,env:HA_WEBHOOK"`
}

type ConfigHTTP struct {
	Host             string  `arg:"--http-host,env:HTTP_HOST" default:"0.0.0.0"`
	Port             uint    `arg:"--http-port,env:HTTP_PORT" default:"8080"`
	Tls              bool    `arg:"--http-tls,env:HTTP_TLS" default:"true"`
	TlsKey           *string `arg:"--http-tls-key,env:HTTP_TLS_KEY"`
	TlsCert          *string `arg:"--http-tls-cert,env:HTTP_TLS_CERT"`
	PathGetLiveness  string  `arg:"env:HTTP_PATH_LIVENESS" default:"/healthz"`
	PathGetReadiness string  `arg:"env:HTTP_PATH_READINESS" default:"/readyz"`
	StaticPath       *string `arg:"--http-static,env:HTTP_STATIC"`
}

type ConfigVideoOutputStream struct {
	Source    string      `arg:"--video-out-src,env:VIDEO_OUT_SRC" default:"v4l2src"`
	Device    string      `arg:"--video-out-device,env:VIDEO_OUT_DEVICE" default:"/dev/video0"`
	Codec     StreamCodec `arg:"--video-out-codec,env:VIDEO_OUT_CODEC" default:"vp8"`
	Height    uint        `arg:"--video-out-height,env:VIDEO_OUT_HEIGHT" default:"480"`
	Width     uint        `arg:"--video-out-width,env:VIDEO_OUT_WIDTH" default:"640"`
	USE_QUEUE bool        `arg:"--video-out-queue,env:VIDEO_OUT_QUEUE" default:"false"`
}

type ConfigAudioOutputStream struct {
	Source     string      `arg:"--audio-out-src,env:AUDIO_OUT_SRC" default:"alsasrc"`
	DeviceName string      `arg:"--audio-out-device-name,env:AUDIO_OUT_DEVICE" default:"default"`
	Device     *string     `arg:"--audio-out-device,env:AUDIO_OUT_DEVICE"`
	Codec      StreamCodec `arg:"--audio-out-codec,env:AUDIO_OUT_CODEC" default:"opus"`
	Channels   uint        `arg:"--audio-out-channels,env:AUDIO_OUT_CHANNELS" default:"1"`
	USE_QUEUE  bool        `arg:"--audio-out-queue,env:AUDIO_OUT_QUEUE" default:"false"`
}

type ConfigAudioInputStream struct {
	Sink       string  `arg:"--audio-in-src,env:AUDIO_IN_SINK" default:"alsasink"`
	DeviceName string  `arg:"--audio-in-device-name,env:AUDIO_IN_DEVICE" default:"default"`
	Device     *string `arg:"--audio-in-device,env:AUDIO_IN_DEVICE"`
	Codec      string  `arg:"--audio-in-codec,env:AUDIO_IN_CODEC" default:"opus"`
	Channels   uint    `arg:"--audio-in-channels,env:AUDIO_IN_CHANNELS" default:"1"`
}

func (c *Config) Stream() *ConfigStream {
	return &ConfigStream{
		VideoOut: c.VideoOut,
		AudioOut: c.AudioOut,
		AudioIn:  c.AudioIn,
	}
}

func (c *ConfigHTTP) Address() string {
	return net.JoinHostPort(c.Host, fmt.Sprint(c.Port))
}

func (c *Config) MustParse() {
	arg.MustParse(c)
}

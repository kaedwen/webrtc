package common

import (
	"fmt"
	"net"

	"github.com/alexflint/go-arg"
)

type Config struct {
	VideoSrc  ConfigVideoSourceStream
	AudioSrc  ConfigAudioSourceStream
	AudioSink ConfigAudioSinkStream
	Logging   ConfigLogging
	Ring      ConfigRing
	Http      ConfigHTTP
}

type ConfigStream struct {
	VideoSrc  ConfigVideoSourceStream // video src for webrtc send
	AudioSrc  ConfigAudioSourceStream // audio src for webrtc send
	AudioSink ConfigAudioSinkStream   // audio sink for webrtc receive
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
	NoIPv6               bool    `arg:"--disable-ipv6,env:NO_IPV6" default:"false"`
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

type ConfigVideoSourceStream struct {
	Source string      `arg:"--video-src,env:VIDEO_SRC" default:"v4l2src"`
	Device string      `arg:"--video-src-device,env:VIDEO_SRC_DEVICE" default:"/dev/video0"`
	Codec  StreamCodec `arg:"--video-src-codec,env:VIDEO_SRC_CODEC" default:"vp8"`
	Height uint        `arg:"--video-src-height,env:VIDEO_SRC_HEIGHT" default:"480"`
	Width  uint        `arg:"--video-src-width,env:VIDEO_SRC_WIDTH" default:"640"`
	Queue  bool        `arg:"--video-src-queue,env:VIDEO_SRC_QUEUE" default:"false"`
}

type ConfigAudioSourceStream struct {
	Source     string      `arg:"--audio-src,env:AUDIO_SRC" default:"alsasrc"`
	DeviceName string      `arg:"--audio-src-device-name,env:AUDIO_SRC_DEVICE_NAME" default:"default"`
	Device     *string     `arg:"--audio-src-device,env:AUDIO_SRC_DEVICE"`
	Codec      StreamCodec `arg:"--audio-src-codec,env:AUDIO_SRC_CODEC" default:"opus"`
	Channels   uint        `arg:"--audio-src-channels,env:AUDIO_SRC_CHANNELS" default:"1"`
	Queue      bool        `arg:"--audio-src-queue,env:AUDIO_SRC_QUEUE" default:"false"`
}

type ConfigAudioSinkStream struct {
	Sink       string  `arg:"--audio-sink,env:AUDIO_SINK" default:"alsasink"`
	DeviceName string  `arg:"--audio-sink-device-name,env:AUDIO_SINK_DEVICE_NAME" default:"default"`
	Device     *string `arg:"--audio-sink-device,env:AUDIO_SINK_DEVICE"`
	Codec      string  `arg:"--audio-sink-codec,env:AUDIO_SINK_CODEC" default:"opus"`
	Channels   uint    `arg:"--audio-sink-channels,env:AUDIO_SINK_CHANNELS" default:"1"`
	Queue      bool    `arg:"--audio-sink-queue,env:AUDIO_SINK_QUEUE" default:"false"`
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

func (c *Config) MustParse() {
	arg.MustParse(&c.VideoSrc)
	arg.MustParse(&c.AudioSrc)
	arg.MustParse(&c.AudioSink)
	arg.MustParse(&c.Logging)
	arg.MustParse(&c.Http)
	arg.MustParse(&c.Ring)
}

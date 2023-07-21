package common

import (
	"fmt"
	"strings"

	"github.com/pion/webrtc/v3"
)

type StreamCodec string

const (
	// video codecs
	H264 StreamCodec = "H264"
	VP8  StreamCodec = "VP8"
	VP9  StreamCodec = "VP9"

	// audio codecs
	OPUS StreamCodec = "OPUS"
)

func (c *StreamCodec) UnmarshalText(text []byte) error {
	switch strings.ToUpper(string(text)) {
	case string(H264):
		*c = H264
	case string(VP8):
		*c = VP8
	case string(VP9):
		*c = VP9
	case string(OPUS):
		*c = OPUS
	default:
		return fmt.Errorf("unsupported codec - %s", text)
	}

	return nil
}

func (c StreamCodec) Mime() string {
	switch c {
	case H264:
		return webrtc.MimeTypeH264
	case VP8:
		return webrtc.MimeTypeVP8
	case VP9:
		return webrtc.MimeTypeVP9
	case OPUS:
		return webrtc.MimeTypeOpus
	default:
		return "UNKNOWN"
	}
}

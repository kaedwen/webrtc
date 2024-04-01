package webrtc

import (
	"context"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/go-gst/go-gst/gst"
	"github.com/kaedwen/webrtc/pkg/common"
	"github.com/kaedwen/webrtc/pkg/server"
	"github.com/kaedwen/webrtc/pkg/streamer"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
	"go.uber.org/zap"
)

const STUN_SERVER = "stun:stun.l.google.com:19302"

type WebrtcHandler struct {
	lg            *zap.Logger
	mu            *sync.Mutex
	cfg           *common.ConfigStream
	audioPipeline *gst.Pipeline
	videoPipeline *gst.Pipeline
	peerHandles   map[string]*PeerHandle
}

type PeerHandle struct {
	audioTrack *webrtc.TrackLocalStaticSample
	videoTrack *webrtc.TrackLocalStaticSample
}

func NewWebrtcHandler(ctx context.Context, lg *zap.Logger, cfg *common.ConfigStream, ch <-chan *server.SignalingHandle) error {
	wh := WebrtcHandler{
		lg:          lg,
		cfg:         cfg,
		mu:          &sync.Mutex{},
		peerHandles: make(map[string]*PeerHandle, 0),
	}

	err := wh.handleAudioSamples(ctx, &cfg.AudioSrc)
	if err != nil {
		return err
	}

	err = wh.handleVideoSamples(ctx, &cfg.VideoSrc)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case sh := <-ch:
				wh.lg.Info("new connection", zap.String("id", sh.Id))

				// create new peer handle
				wh.createPeerHandle(ctx, sh)

				// make sure the pipelines are running
				wh.startPipelines()

			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (wh *WebrtcHandler) startPipelines() {

	if wh.audioPipeline.GetCurrentState() != gst.StatePlaying {
		err := wh.audioPipeline.SetState(gst.StatePlaying)
		if err != nil {
			wh.lg.Fatal("failed to start audio pipeline", zap.Error(err))
		}
		wh.lg.Info("started audio pipeline")
	}

	if wh.videoPipeline.GetCurrentState() != gst.StatePlaying {
		err := wh.videoPipeline.SetState(gst.StatePlaying)
		if err != nil {
			wh.lg.Fatal("failed to start video pipeline", zap.Error(err))
		}
		wh.lg.Info("started video pipeline")
	}

}

func (wh *WebrtcHandler) stopPipelines() {
	var err error

	err = wh.audioPipeline.SetState(gst.StateNull)
	if err != nil {
		wh.lg.Fatal("failed to pause audio pipeline", zap.Error(err))
	}

	err = wh.videoPipeline.SetState(gst.StateNull)
	if err != nil {
		wh.lg.Fatal("failed to pause video pipeline", zap.Error(err))
	}

}

func (wh *WebrtcHandler) handleAudioSamples(ctx context.Context, cfg *common.ConfigAudioSourceStream) error {
	properties := map[string]interface{}{}
	if cfg.Source == "alsasrc" || cfg.Source == "pulsesrc" {
		if cfg.Device != nil {
			properties["device"] = *cfg.Device
		} else {
			properties["device"] = cfg.DeviceName
		}
	}

	src := streamer.StreamElement{
		Kind:       cfg.Source,
		Properties: properties,
		Caps:       streamer.NewCapsBuilder("audio/x-raw").Channels(int(cfg.Channels)).Rate(48000),
		Queue:      cfg.Queue,
		Codec:      cfg.Codec,
	}

	var err error
	var audioCh <-chan media.Sample
	wh.audioPipeline, audioCh, err = streamer.CreateAudioPipelineSink(src, wh.lg)
	if err != nil {
		return err
	}

	streamer.LoopBus(wh.lg.With(zap.String("sub-context", "audio")), wh.audioPipeline)

	go func() {
		wh.lg.Info("wait for audio sample")
		for {
			select {
			case data := <-audioCh:
				for id, ph := range wh.peerHandles {
					err := ph.audioTrack.WriteSample(data)
					if err != nil {
						wh.lg.Error("failed to write audio sample", zap.String("id", id), zap.Error(err))
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (wh *WebrtcHandler) handleVideoSamples(ctx context.Context, cfg *common.ConfigVideoSourceStream) error {
	// src := streamer.StreamElement{
	// 	Kind: cfg.Source,
	// 	Properties: map[string]interface{}{
	// 		"device": cfg.Device,
	// 	},
	// 	Caps:  streamer.NewCapsBuilder("video/x-raw").Format("YUY2").Width(int(cfg.Width)).Height(int(cfg.Height)),
	// 	Queue: cfg.Queue,
	// 	Codec: cfg.Codec,
	// }

	var err error
	var videoCh <-chan media.Sample
	wh.videoPipeline, videoCh, err = streamer.CreateVideoPipelineSinkWithLaunch(wh.lg)
	if err != nil {
		return err
	}

	streamer.LoopBus(wh.lg.With(zap.String("sub-context", "video")), wh.videoPipeline)

	go func() {
		wh.lg.Info("wait for video sample")
		for {
			select {
			case data := <-videoCh:
				for id, ph := range wh.peerHandles {
					err := ph.videoTrack.WriteSample(data)
					if err != nil {
						wh.lg.Error("failed to write video sample", zap.String("id", id), zap.Error(err))
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (wh *WebrtcHandler) createPeerHandle(rctx context.Context, sh *server.SignalingHandle) error {
	wh.mu.Lock()
	defer wh.mu.Unlock()

	// Prepare the configuration
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{{URLs: []string{STUN_SERVER}}},
	}

	// Create a new RTCPeerConnection
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		return err
	}

	// create a context for this handle
	hctx, hcancel := context.WithCancel(rctx)

	// sent the candidate out when where is one
	peerConnection.OnICECandidate(func(i *webrtc.ICECandidate) {
		if i != nil {
			sh.Trcv <- server.NewIceCandidateMessage(*i)
		}
	})

	// Set a handler for when a new remote track starts, this handler creates a gstreamer pipeline
	// for the given codec
	peerConnection.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		wh.lg.Info("received track", zap.String("kind", track.Kind().String()), zap.String("codec", track.Codec().MimeType))
		cfg := wh.cfg.AudioSink

		if track.Codec().MimeType != "audio/opus" {
			wh.lg.Error("mimetype not supported", zap.String("mime", track.Codec().MimeType))
			return
		}

		properties := map[string]interface{}{}
		if cfg.Sink == "alsasink" || cfg.Sink == "pulsesink" {
			if cfg.Device != nil {
				properties["device"] = *cfg.Device
			} else {
				properties["device"] = cfg.DeviceName
			}
		}

		// Send a PLI on an interval so that the publisher is pushing a keyframe every rtcpPLIInterval
		go func() {
			ticker := time.NewTicker(time.Second * 3)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					if err := peerConnection.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: uint32(track.SSRC())}}); err != nil {
						if errors.Is(err, io.ErrClosedPipe) {
							return
						}

						wh.lg.Error("failed to send PLI", zap.Error(err))
					}
				case <-hctx.Done():
					return
				}
			}
		}()

		pipeline, err := streamer.CreateAudioPipelineSrc(streamer.StreamElement{
			Kind:       cfg.Sink,
			Properties: properties,
			Queue:      cfg.Queue,
		})
		if err != nil {
			wh.lg.Error("failed to create src pipeline", zap.Error(err))
			return
		}

		err = pipeline.Start()
		if err != nil {
			wh.lg.Error("failed to start pipeline", zap.Error(err))
			return
		}

		buf := make([]byte, 1400)
		for {
			n, _, readErr := track.Read(buf)
			if readErr != nil {
				if errors.Is(readErr, io.EOF) {
					return
				}

				wh.lg.Error("read failed", zap.Error(err))
				continue
			}

			if n > 0 {
				err := pipeline.Push(buf[:n])
				if err != nil {
					wh.lg.Error("push failed", zap.Error(err))
				}
			}
		}
	})

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		wh.lg.Info("connection-state has changed", zap.String("state", connectionState.String()))

		if connectionState == webrtc.ICEConnectionStateDisconnected {
			peerConnection.Close()
			hcancel()

			// remove this handle
			wh.mu.Lock()

			delete(wh.peerHandles, sh.Id)

			if len(wh.peerHandles) == 0 {
				wh.stopPipelines()
			}

			wh.mu.Unlock()
		}
	})

	// Create a audio track
	audioTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: wh.cfg.AudioSrc.Codec.Mime()}, "audio", "pion1")
	if err != nil {
		return err
	}
	_, err = peerConnection.AddTrack(audioTrack)
	if err != nil {
		return err
	}

	// Create a video track
	videoTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: wh.cfg.VideoSrc.Codec.Mime()}, "video", "pion2")
	if err != nil {
		return err
	}
	_, err = peerConnection.AddTrack(videoTrack)
	if err != nil {
		return err
	}

	onOfferReceived := func(offer webrtc.SessionDescription) error {

		// Set the remote SessionDescription
		err = peerConnection.SetRemoteDescription(offer)
		if err != nil {
			return err
		}

		// Create an answer
		answer, err := peerConnection.CreateAnswer(nil)
		if err != nil {
			return err
		}

		// Create channel that is blocked until ICE Gathering is complete
		//gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

		// Sets the LocalDescription, and starts our UDP listeners
		err = peerConnection.SetLocalDescription(answer)
		if err != nil {
			return err
		}

		//<-gatherComplete

		// Send the answer
		sh.Trcv <- server.NewAnswerMessage(peerConnection.LocalDescription())

		return nil
	}

	hndl := PeerHandle{audioTrack, videoTrack}

	// add handle to list
	wh.peerHandles[sh.Id] = &hndl

	go func() {
		for {
			select {
			case m, ok := <-sh.Recv:
				// return when channel is closed
				if !ok {
					wh.mu.Lock()
					defer wh.mu.Unlock()

					// remove this handle
					delete(wh.peerHandles, sh.Id)

					// when there is no left over pause the pipelines
					if len(wh.peerHandles) == 0 {
						wh.stopPipelines()
					}

					return
				}

				switch true {
				case m.IsIceCandidateMessage():
					pm, err := m.ToIceCandidateMessage()
					if err != nil {
						wh.lg.Error("failed to parse candidate", zap.Error(err))
						continue
					}

					wh.lg.Info("received new ice candidate", zap.Any("candidate", pm.Candidate))

					err = peerConnection.AddICECandidate(pm.Candidate)
					if err != nil {
						wh.lg.Error("failed to add ice candidate", zap.Error(err))
					}
				case m.IsAnswerMessage():
					wh.lg.Info("received answer", zap.Any("data", m.Data))
				case m.IsOfferMessage():
					pm, err := m.ToOfferMessage()
					if err != nil {
						wh.lg.Error("failed to parse offer", zap.Error(err))
						continue
					}

					wh.lg.Info("received offer", zap.Any("data", pm.Offer))

					err = onOfferReceived(pm.Offer)
					if err != nil {
						wh.lg.Error("failed to handle offer", zap.Error(err))
					}
				}
			case <-hctx.Done():
				return
			}
		}
	}()

	return nil
}

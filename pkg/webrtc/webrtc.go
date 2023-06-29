package webrtc

import (
	"context"
	"sync"

	"gitea.heinrich.blue/PHI/webrtc-gst/pkg/server"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
	"github.com/tinyzimmer/go-gst/gst"
	"go.uber.org/zap"
)

const STUN_SERVER = "stun:stun.l.google.com:19302"

type WebrtcHandler struct {
	lg            *zap.Logger
	mu            *sync.Mutex
	audioPipeline *gst.Pipeline
	videoPipeline *gst.Pipeline
	peerHandles   map[string]*PeerHandle
}

type PeerHandle struct {
	audioTrack *webrtc.TrackLocalStaticSample
	videoTrack *webrtc.TrackLocalStaticSample
}

func NewWebrtcHandler(ctx context.Context, lg *zap.Logger, ch <-chan *server.SignalingHandle) error {
	wh := WebrtcHandler{
		lg:          lg,
		mu:          &sync.Mutex{},
		peerHandles: make(map[string]*PeerHandle, 0),
	}

	err := wh.handleAudioSamples(ctx, "audiotestsrc")
	if err != nil {
		return err
	}

	err = wh.handleVideoSamples(ctx, "videotestsrc")
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

	if wh.audioPipeline.GetState() != gst.StatePlaying {
		err := wh.audioPipeline.SetState(gst.StatePlaying)
		if err != nil {
			wh.lg.Fatal("failed to start audio pipeline", zap.Error(err))
		}
		wh.lg.Info("started audio pipeline")
	}

	if wh.videoPipeline.GetState() != gst.StatePlaying {
		err := wh.videoPipeline.SetState(gst.StatePlaying)
		if err != nil {
			wh.lg.Fatal("failed to start video pipeline", zap.Error(err))
		}
		wh.lg.Info("started video pipeline")
	}

}

func (wh *WebrtcHandler) stopPipelines() {
	var err error

	err = wh.audioPipeline.SetState(gst.StatePaused)
	if err != nil {
		wh.lg.Fatal("failed to pause audio pipeline", zap.Error(err))
	}

	err = wh.videoPipeline.SetState(gst.StatePaused)
	if err != nil {
		wh.lg.Fatal("failed to pause video pipeline", zap.Error(err))
	}

}

func (wh *WebrtcHandler) handleAudioSamples(ctx context.Context, src string) error {
	var err error
	var audioCh <-chan media.Sample
	wh.audioPipeline, audioCh, err = CreateAudioPipeline(src)
	if err != nil {
		return err
	}

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

func (wh *WebrtcHandler) handleVideoSamples(ctx context.Context, src string) error {
	var err error
	var videoCh <-chan media.Sample
	wh.videoPipeline, videoCh, err = CreateVideoPipeline(src)
	if err != nil {
		return err
	}

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

func (wh *WebrtcHandler) createPeerHandle(ctx context.Context, sh *server.SignalingHandle) error {
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

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		wh.lg.Info("connection-state has changed", zap.String("state", connectionState.String()))

		if connectionState == webrtc.ICEConnectionStateDisconnected {
			peerConnection.Close()

			// remove this handle
			wh.mu.Lock()
			delete(wh.peerHandles, sh.Id)
			wh.mu.Unlock()
		}
	})

	// Create a audio track
	audioTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "audio/opus"}, "audio", "pion1")
	if err != nil {
		return err
	}
	_, err = peerConnection.AddTrack(audioTrack)
	if err != nil {
		return err
	}

	// Create a video track
	videoTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: "video/vp8"}, "video", "pion2")
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
		gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

		// Sets the LocalDescription, and starts our UDP listeners
		err = peerConnection.SetLocalDescription(answer)
		if err != nil {
			return err
		}

		<-gatherComplete

		// Send the answer
		sh.Trcv <- *peerConnection.LocalDescription()

		return nil
	}

	hndl := PeerHandle{
		audioTrack: audioTrack,
		videoTrack: videoTrack,
	}

	go func() {
		for {
			select {
			case offer, ok := <-sh.Recv:
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

				err := onOfferReceived(offer)
				if err != nil {
					wh.lg.Error("failed to handle offer", zap.Error(err))
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// add handle to list
	wh.peerHandles[sh.Id] = &hndl

	return nil
}

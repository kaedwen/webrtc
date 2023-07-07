package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kaedwen/webrtc/pkg/common"
	"github.com/kaedwen/webrtc/static"
	"github.com/pion/webrtc/v3"
	"go.uber.org/zap"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type SignalingHandle struct {
	Id   string
	Recv chan webrtc.SessionDescription
	Trcv chan webrtc.SessionDescription
}

type HttpServer struct {
	http.Server
	lg   *zap.Logger
	Hndl chan *SignalingHandle
}

func NewSignalingHandle(id string) SignalingHandle {
	return SignalingHandle{
		Id:   id,
		Recv: make(chan webrtc.SessionDescription, 10),
		Trcv: make(chan webrtc.SessionDescription, 10),
	}
}

func NewHttpServer(lg *zap.Logger, cfg *common.Config) *HttpServer {
	h := HttpServer{
		Hndl: make(chan *SignalingHandle, 10),
		lg:   lg,
	}

	engine := gin.Default()
	engine.GET("/signaling/:id", h.signalingHandler)

	// static handler
	static.SetupHandler(engine, cfg.Http)

	// set out handler
	h.Handler = engine

	return &h
}

func (h *HttpServer) ListenAndServe(ctx context.Context, addr string) {
	go func() {

		// set the configured address
		h.Addr = addr

		// and listen
		if err := h.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			h.lg.Fatal("listen failed", zap.Error(err))
		}

	}()
}

func (h *HttpServer) TearDown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.Server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %s", err)
	}

	return nil
}

func (h *HttpServer) signalingHandler(c *gin.Context) {
	id := c.Param("id")

	conn, err := websocket.Accept(c.Writer, c.Request, nil)
	if err != nil {
		h.lg.Error("failed to upgrade websocket", zap.Error(err))
		c.Status(http.StatusBadRequest)
		return
	}

	// tear down when going down
	defer conn.Close(websocket.StatusInternalError, "the sky is falling")

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	hndl := NewSignalingHandle(id)

	// push hndl outside
	h.Hndl <- &hndl

	// loop read
	go func() {
		for {
			// wait and read/parse message
			var sdp webrtc.SessionDescription
			err := wsjson.Read(ctx, conn, &sdp)
			if err != nil {
				status := websocket.CloseStatus(err)
				if status == websocket.StatusNormalClosure || status == websocket.StatusGoingAway {

					// close channels to signal closure
					close(hndl.Recv)
					close(hndl.Trcv)

					h.lg.Info("socket closed")
					break
				}

				h.lg.Error("failed to decode message", zap.Error(err))
				break
			}

			h.lg.Info("received message", zap.String("type", sdp.Type.String()))

			// forward in channel
			hndl.Recv <- sdp
		}

		cancel()
	}()

	// loop write
	for sdp := range hndl.Trcv {
		err := wsjson.Write(ctx, conn, sdp)
		if err != nil {
			status := websocket.CloseStatus(err)
			if status == websocket.StatusNormalClosure || status == websocket.StatusGoingAway {
				h.lg.Info("socket closed")
				break
			}

			h.lg.Error("failed to encode message", zap.Error(err))
			break
		}

		h.lg.Info("tranceived message", zap.String("type", sdp.Type.String()))
	}

}

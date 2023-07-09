package server

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"log"
	"math/big"
	"net"
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
	cfg  *common.ConfigHTTP
	rs   *http.Server
	Hndl chan *SignalingHandle
}

func NewSignalingHandle(id string) SignalingHandle {
	return SignalingHandle{
		Id:   id,
		Recv: make(chan webrtc.SessionDescription, 10),
		Trcv: make(chan webrtc.SessionDescription, 10),
	}
}

func NewHttpServer(lg *zap.Logger, cfg *common.ConfigHTTP) *HttpServer {
	h := HttpServer{
		Hndl: make(chan *SignalingHandle, 10),
		cfg:  cfg,
		lg:   lg,
	}

	engine := gin.Default()
	engine.GET("/signaling/:id", h.signalingHandler)

	// static handler
	static.SetupHandler(engine, cfg)

	// set out handler
	h.Handler = engine

	return &h
}

func (h *HttpServer) ListenAndServe(ctx context.Context) error {
	if h.cfg.Tls {
		// set the configured address
		h.Addr = h.cfg.AddressTls()

		if h.cfg.TlsCert != nil && h.cfg.TlsKey != nil {
			cert, err := tls.LoadX509KeyPair(*h.cfg.TlsCert, *h.cfg.TlsKey)
			if err != nil {
				return err
			}

			h.Server.TLSConfig = &tls.Config{
				Certificates: []tls.Certificate{cert},
			}
		} else {
			h.lg.Info("creating self signed certificate")
			cert, err := generateSelfSigned()
			if err != nil {
				return err
			}

			h.Server.TLSConfig = &tls.Config{
				Certificates: []tls.Certificate{cert},
			}
		}

		go func() {
			// and listen
			if err := h.Server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
				h.lg.Fatal("listen failed", zap.Error(err))
			}
		}()

		h.rs = runRedirect(h.cfg)
	} else {
		// set the configured address
		h.Addr = h.cfg.Address()

		go func() {
			// and listen
			if err := h.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				h.lg.Fatal("listen failed", zap.Error(err))
			}
		}()
	}

	return nil
}

func runRedirect(cfg *common.ConfigHTTP) *http.Server {
	engine := gin.Default()

	engine.NoRoute(func(c *gin.Context) {
		u := c.Request.URL

		host, _, _ := net.SplitHostPort(c.Request.Host)
		u.Host = net.JoinHostPort(host, fmt.Sprint(cfg.PortTls))
		u.Scheme = "https"

		c.Redirect(302, u.String())
	})

	go func() {
		engine.Run(net.JoinHostPort(cfg.Host, fmt.Sprint(cfg.Port)))
	}()

	return &http.Server{
		Addr:    cfg.Address(),
		Handler: engine,
	}
}

func generateSelfSigned() (tls.Certificate, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatal(err)
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to generate serial number: %v", err)
	}

	template := x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               pkix.Name{Organization: []string{"PHI"}},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 180),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	cert, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return tls.Certificate{}, err
	}

	return tls.Certificate{
		PrivateKey:  priv,
		Certificate: [][]byte{cert},
	}, nil
}

func (h *HttpServer) TearDown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if h.rs != nil {
		if err := h.rs.Shutdown(ctx); err != nil {
			return fmt.Errorf("redirect server forced to shutdown: %s", err)
		}
	}

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

//go:build embed

package static

import (
	"embed"
	"io"
	"io/fs"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/kaedwen/webrtc/pkg/common"
)

//go:embed dist/*
var f embed.FS

type staticSourceFile struct {
	target fs.File
}

func (ss *staticSourceFile) Reader() io.ReadCloser {
	return ss.target
}

func (ss *staticSourceFile) Stat() (fs.FileInfo, error) {
	return ss.target.Stat()
}

type staticSource struct {
	base string
}

func (s staticSource) Open(p string) (common.StaticSourceFile, error) {
	t, err := f.Open(path.Join(s.base, p))
	if err != nil {
		return nil, err
	}

	return &staticSourceFile{t}, nil
}

func SetupHandler(e *gin.Engine, cfg *common.Config) {
	if cfg.Ring.JinglePath != nil {
		e.GET(cfg.Ring.JingleName, func(c *gin.Context) {
			c.File(path.Join(*cfg.Ring.JinglePath, cfg.Ring.JingleName))
		})
	}

	handler := common.NewStaticHandler(staticSource{"dist"}, "index.html")
	e.NoRoute(func(c *gin.Context) {
		encoding := c.Request.Header.Get("Accept-Encoding")

		p := c.Request.URL.Path
		if i, err := handler.Get(p, encoding); err == nil {
			c.DataFromReader(http.StatusOK, i.Size, i.Mime, i.Reader, i.ExtraHeaders)
			return
		}

		c.Status(http.StatusNotFound)
	})
}

//go:build !embed

package static

import (
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/kaedwen/webrtc/pkg/common"
)

type staticSourceFile struct {
	target *os.File
}

func (ss *staticSourceFile) Reader() io.ReadCloser {
	return ss.target
}

func (ss *staticSourceFile) Stat() (fs.FileInfo, error) {
	return ss.target.Stat()
}

type staticDiskSource struct {
	base string
}

func (s staticDiskSource) Open(p string) (common.StaticSourceFile, error) {
	t, err := os.Open(path.Join(s.base, p))
	if err != nil {
		return nil, err
	}

	return &staticSourceFile{t}, nil
}

func SetupHandler(e *gin.Engine, cfg *common.Config) {
	if cfg.Http.StaticPath != nil {
		handler := common.NewStaticHandler(staticDiskSource{cfg.Http.StaticPath.String()}, "index.html")
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
}

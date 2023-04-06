package compress

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/go-labx/lightning"
	"github.com/klauspost/compress/zstd"
)

type config struct {
	writer *bytes.Buffer
}

type Options func(*config)

func Default() lightning.Middleware {
	return New()
}

func New(options ...Options) lightning.Middleware {
	cfg := &config{}

	for _, option := range options {
		option(cfg)
	}

	buf := &bytes.Buffer{}
	brWriter := brotli.NewWriter(cfg.writer)
	deflateWriter, _ := flate.NewWriter(cfg.writer, 5)
	gzipWriter := gzip.NewWriter(cfg.writer)
	zstdWriter, _ := zstd.NewWriter(cfg.writer)

	return func(ctx *lightning.Context) {
		ctx.Next()

		acceptEncoding := ctx.Header("Accept-Encoding")
		body := ctx.Body()

		switch {
		case strings.Contains(acceptEncoding, "br"):
			brWriter.Reset(buf)
			buf.Reset()
			_, err := brWriter.Write(body)
			if err != nil {
				return
			}
			err = brWriter.Flush()
			if err != nil {
				return
			}

			ctx.SetHeader("Content-Encoding", "br")
			ctx.SetBody(buf.Bytes())
		case strings.Contains(acceptEncoding, "deflate"):
			deflateWriter.Reset(buf)
			buf.Reset()

			_, err := deflateWriter.Write(body)
			if err != nil {
				return
			}
			err = deflateWriter.Flush()
			if err != nil {
				return
			}

			ctx.SetHeader("Content-Encoding", "deflate")
			ctx.SetBody(buf.Bytes())
		case strings.Contains(acceptEncoding, "gzip"):
			gzipWriter.Reset(buf)
			buf.Reset()

			_, err := gzipWriter.Write(body)
			if err != nil {
				return
			}
			err = gzipWriter.Flush()
			if err != nil {
				return
			}

			ctx.SetHeader("Content-Encoding", "gzip")
			ctx.SetBody(buf.Bytes())
		case strings.Contains(acceptEncoding, "zstd"):
			zstdWriter.Reset(buf)
			buf.Reset()

			_, err := zstdWriter.Write(body)
			if err != nil {
				return
			}
			err = zstdWriter.Flush()
			if err != nil {
				return
			}

			ctx.SetHeader("Content-Encoding", "zstd")
			ctx.SetBody(buf.Bytes())
		}
	}
}

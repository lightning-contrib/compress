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

type Compressor interface {
	Compress([]byte) ([]byte, error)
}

type BrotliCompression struct{}

func (c *BrotliCompression) Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := brotli.NewWriter(&buf)
	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}
	err = writer.Flush()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type DeflateCompression struct{}

func (c *DeflateCompression) Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer, err := flate.NewWriter(&buf, 5)
	if err != nil {
		return nil, err
	}
	_, err = writer.Write(data)
	if err != nil {
		return nil, err
	}
	err = writer.Flush()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type GzipCompression struct{}

func (c *GzipCompression) Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}
	err = writer.Flush()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type ZstdCompression struct{}

func (c *ZstdCompression) Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer, err := zstd.NewWriter(&buf)
	if err != nil {
		return nil, err
	}
	_, err = writer.Write(data)
	if err != nil {
		return nil, err
	}
	err = writer.Flush()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// br [0,11]
// deflate / gzip [-2,9]
// zstd [0,5]
type config struct {
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

	return func(ctx *lightning.Context) {
		ctx.Next()

		acceptEncoding := ctx.Header("Accept-Encoding")
		body := ctx.Body()

		var encoding string
		var compressor Compressor

		switch {
		case strings.Contains(acceptEncoding, "br"):
			compressor = &BrotliCompression{}
			encoding = "br"
		case strings.Contains(acceptEncoding, "deflate"):
			compressor = &DeflateCompression{}
			encoding = "deflate"
		case strings.Contains(acceptEncoding, "gzip"):
			compressor = &GzipCompression{}
			encoding = "gzip"
		case strings.Contains(acceptEncoding, "zstd"):
			compressor = &ZstdCompression{}
			encoding = "zstd"
		}

		compressed, err := compressor.Compress(body)
		if err != nil {
			return
		}
		ctx.SetHeader("Content-Encoding", encoding)
		ctx.SetBody(compressed)
	}
}

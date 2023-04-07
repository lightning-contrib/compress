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

const (
	EncodingBrotli  = "br"
	EncodingDeflate = "deflate"
	EncodingGzip    = "gzip"
	EncodingZstd    = "zstd"
)

// Compressor is an interface that defines the Compress method
type Compressor interface {
	Compress([]byte) ([]byte, error)
}

type BrotliCompression struct{}

// Compress compresses the given data using Brotli compression algorithm
func (c *BrotliCompression) Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := brotli.NewWriterLevel(&buf, brotli.DefaultCompression)
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

// Compress compresses the given data using Deflate compression algorithm
func (c *DeflateCompression) Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer, err := flate.NewWriter(&buf, flate.DefaultCompression)
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

// Compress compresses the given data using Gzip compression algorithm
func (c *GzipCompression) Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer, err := gzip.NewWriterLevel(&buf, gzip.DefaultCompression)
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

type ZstdCompression struct{}

// Compress compresses the given data using Zstd compression algorithm
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

// Options is a function that takes a pointer to a config struct
type Options func(*config)

// Default returns a lightning middleware with default options
func Default() lightning.Middleware {
	return New()
}

// New returns a lightning middleware with the given options
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
		case strings.Contains(acceptEncoding, EncodingBrotli):
			compressor = &BrotliCompression{}
			encoding = EncodingBrotli
		case strings.Contains(acceptEncoding, EncodingDeflate):
			compressor = &DeflateCompression{}
			encoding = EncodingDeflate
		case strings.Contains(acceptEncoding, EncodingGzip):
			compressor = &GzipCompression{}
			encoding = EncodingGzip
		case strings.Contains(acceptEncoding, EncodingZstd):
			compressor = &ZstdCompression{}
			encoding = EncodingZstd
		}

		compressed, err := compressor.Compress(body)
		if err != nil {
			return
		}
		ctx.SetHeader("Content-Encoding", encoding)
		ctx.SetBody(compressed)
	}
}

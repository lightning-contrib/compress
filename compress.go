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
	Compress(data []byte, level int) ([]byte, error)
}

type BrotliCompression struct{}

// Compress compresses the given data using Brotli compression algorithm
func (c *BrotliCompression) Compress(data []byte, level int) ([]byte, error) {
	var buf bytes.Buffer
	writer := brotli.NewWriterLevel(&buf, level)
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
func (c *DeflateCompression) Compress(data []byte, level int) ([]byte, error) {
	var buf bytes.Buffer
	writer, err := flate.NewWriter(&buf, level)
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
func (c *GzipCompression) Compress(data []byte, level int) ([]byte, error) {
	var buf bytes.Buffer
	writer, err := gzip.NewWriterLevel(&buf, level)
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
func (c *ZstdCompression) Compress(data []byte, _ int) ([]byte, error) {
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

type config struct {
	// BrotliCompressionLevel is the compression level for Brotli compression algorithm
	BrotliCompressionLevel int
	// DeflateCompressionLevel is the compression level for Deflate compression algorithm
	DeflateCompressionLevel int
	// GzipCompressionLevel is the compression level for Gzip compression algorithm
	GzipCompressionLevel int
}

// Options is a function that takes a pointer to a config struct
type Options func(*config)

// WithBrotliCompressionLevel sets the compression level for Brotli compression algorithm
func WithBrotliCompressionLevel(level int) Options {
	return func(c *config) {
		c.BrotliCompressionLevel = level
	}
}

// WithDeflateCompressionLevel sets the compression level for Deflate compression algorithm
func WithDeflateCompressionLevel(level int) Options {
	return func(c *config) {
		c.DeflateCompressionLevel = level
	}
}

// WithGzipCompressionLevel sets the compression level for Gzip compression algorithm
func WithGzipCompressionLevel(level int) Options {
	return func(c *config) {
		c.GzipCompressionLevel = level
	}
}

// Default returns a lightning middleware with default options
func Default() lightning.Middleware {
	return New()
}

// New returns a lightning middleware with the given options
func New(options ...Options) lightning.Middleware {
	cfg := &config{
		BrotliCompressionLevel:  brotli.DefaultCompression,
		DeflateCompressionLevel: flate.DefaultCompression,
		GzipCompressionLevel:    gzip.DefaultCompression,
	}

	for _, option := range options {
		option(cfg)
	}

	return func(ctx *lightning.Context) {
		ctx.Next()

		acceptEncoding := ctx.Header("Accept-Encoding")
		var encoding string
		var compressor Compressor
		var level int

		switch {
		case strings.Contains(acceptEncoding, EncodingBrotli):
			compressor = &BrotliCompression{}
			encoding = EncodingBrotli
			level = cfg.BrotliCompressionLevel
		case strings.Contains(acceptEncoding, EncodingDeflate):
			compressor = &DeflateCompression{}
			encoding = EncodingDeflate
			level = cfg.DeflateCompressionLevel
		case strings.Contains(acceptEncoding, EncodingGzip):
			compressor = &GzipCompression{}
			encoding = EncodingGzip
			level = cfg.GzipCompressionLevel
		case strings.Contains(acceptEncoding, EncodingZstd):
			compressor = &ZstdCompression{}
			encoding = EncodingZstd
		}

		if compressor == nil {
			return
		}
		compressed, err := compressor.Compress(ctx.Body(), level)
		if err != nil {
			return
		}
		ctx.SetHeader("Content-Encoding", encoding)
		ctx.SetBody(compressed)
	}
}

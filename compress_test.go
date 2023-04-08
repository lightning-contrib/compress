package compress

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-labx/lightning"
)

func TestDefault(t *testing.T) {
	middleware := Default()
	if middleware == nil {
		t.Error("Expected middleware to not be nil")
	}
}

func TestNew_1(t *testing.T) {
	// Create a mock request with an empty x-request-id header
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock response recorder
	rr := httptest.NewRecorder()

	ctx, _ := lightning.NewContext(rr, req)

	// Test with default options
	middleware := New()
	ctx.SetBody([]byte("hello world"))

	middleware(ctx)

	contentEncoding := rr.Header().Get("Content-Encoding")

	if contentEncoding != "" {
		t.Errorf("Expected contentEncoding to be an empty string, but got %s", contentEncoding)
	}
}

func TestNew_2(t *testing.T) {
	// Create a mock request with an empty x-request-id header
	req, err := http.NewRequest("GET", "/", nil)
	req.Header.Set("Accept-Encoding", "br")
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock response recorder
	rr := httptest.NewRecorder()

	ctx, _ := lightning.NewContext(rr, req)

	middleware := New(
		WithBrotliCompressionLevel(6),
		WithDeflateCompressionLevel(-1),
		WithGzipCompressionLevel(-1),
	)
	ctx.SetBody([]byte("hello world"))
	middleware(ctx)

	contentEncoding := rr.Header().Get("Content-Encoding")
	if contentEncoding != "br" {
		t.Errorf("Expected contentEncoding to be a br, but got %s", contentEncoding)
	}

	req.Header.Set("Accept-Encoding", "deflate")
	ctx.SetBody([]byte("hello world"))
	middleware(ctx)
	contentEncoding = rr.Header().Get("Content-Encoding")
	if contentEncoding != "deflate" {
		t.Errorf("Expected contentEncoding to be a deflate, but got %s", contentEncoding)
	}

	req.Header.Set("Accept-Encoding", "gzip")
	ctx.SetBody([]byte("hello world"))
	middleware(ctx)
	contentEncoding = rr.Header().Get("Content-Encoding")
	if contentEncoding != "gzip" {
		t.Errorf("Expected contentEncoding to be a gzip, but got %s", contentEncoding)
	}

	req.Header.Set("Accept-Encoding", "zstd")
	ctx.SetBody([]byte("hello world"))
	middleware(ctx)
	contentEncoding = rr.Header().Get("Content-Encoding")
	if contentEncoding != "zstd" {
		t.Errorf("Expected contentEncoding to be an zstd, but got %s", contentEncoding)
	}
}

func TestNew_3(t *testing.T) {
	// Create a mock request with an empty x-request-id header
	req, err := http.NewRequest("GET", "/", nil)
	req.Header.Set("Accept-Encoding", "br")
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock response recorder
	rr := httptest.NewRecorder()

	ctx, _ := lightning.NewContext(rr, req)

	middleware := New()
	req.Header.Set("Accept-Encoding", "xxx")
	ctx.SetBody([]byte("hello world"))
	middleware(ctx)
	contentEncoding := rr.Header().Get("Content-Encoding")
	if contentEncoding != "" {
		t.Errorf("Expected contentEncoding to be an empty string, but got %s", contentEncoding)
	}
}

package main

import (
	"strings"

	"github.com/go-labx/lightning"
	"github.com/lightning-contrib/compress"
)

func main() {
	app := lightning.NewApp()

	app.Use(compress.New(
		compress.WithBrotliCompressionLevel(11),
		compress.WithDeflateCompressionLevel(-1),
		compress.WithGzipCompressionLevel(-1),
	))

	app.Get("/ping", func(ctx *lightning.Context) {
		ctx.Text(200, strings.Repeat("hello world", 100))
	})

	app.Run(":6789")
}

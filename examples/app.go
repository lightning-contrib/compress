package main

import (
	"github.com/go-labx/lightning"
	"github.com/lightning-contrib/compress"
)

func main() {
	app := lightning.DefaultApp()
	app.Use(compress.Default())

	app.Get("/ping", func(ctx *lightning.Context) {
		ctx.Text(200, "hello world")
	})

	app.Run(":6789")
}

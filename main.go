package main

import (
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	cli := &CLI{}
	ctx := kong.Parse(cli)

	err := ctx.Run()
	if err != nil {
		slog.Error("could not execute", slog.String("error", err.Error()))
	}
}

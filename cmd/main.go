package main

import (
	"github.com/folivorra/get_order/internal/config"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(
		slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{
				Level:     slog.LevelDebug,
				AddSource: true,
			},
		),
	)

	cfg := config.NewConfig(logger)

}

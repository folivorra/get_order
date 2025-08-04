package main

import (
	"github.com/folivorra/get_order/internal/config"
	"github.com/joho/godotenv"
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

	if err := godotenv.Load(".env"); err != nil {
		logger.Warn("loading .env file failed, used default values",
			slog.String("error", err.Error()),
		)
	}

	cfg := config.NewConfig(logger)
}

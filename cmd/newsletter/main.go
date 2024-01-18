// Package main is the initial point for the service newsletter.
package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/perebaj/newsletter"
)

// Config is the struct that contains the configuration for the service.
type Config struct {
	LogLevel string
	LogType  string
}

func main() {

	cfg := Config{
		LogLevel: getEnvWithDefault("LOG_LEVEL", "INFO"),
		LogType:  getEnvWithDefault("LOG_TYPE", "json"),
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	if err := setUpLog(cfg); err != nil {
		slog.Error("error setting up log", "error", err)
		signalCh <- syscall.SIGTERM
	}

	sSlice := []string{"http://www.google.com", "www.facebook.com", "www.x.com"}
	jobs := make(chan string, len(sSlice))
	result := make(chan string, len(sSlice))
	for _, s := range sSlice {
		jobs <- s
	}

	go newsletter.Worker(jobs, result, newsletter.GetReferences)

	for i := 0; i < len(sSlice); i++ {
		r := <-result
		if r != "" {
			slog.Info(r)
		}
	}

	<-signalCh
}

// setUpLog initialize the logger.
func setUpLog(cfg Config) error {
	var level slog.Level
	switch cfg.LogLevel {
	case "INFO":
		level = slog.LevelInfo
	case "DEBUG":
		level = slog.LevelDebug
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		return fmt.Errorf("invalid log level: %s", cfg.LogLevel)
	}

	var logger *slog.Logger
	if cfg.LogType == "json" {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		}))
	} else if cfg.LogType == "text" {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		}))
	} else {
		return fmt.Errorf("invalid log type: %s", cfg.LogType)
	}

	slog.SetDefault(logger)
	return nil
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

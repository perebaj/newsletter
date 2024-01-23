// Package main is the initial point for the service newsletter.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/perebaj/newsletter"
	"github.com/perebaj/newsletter/mongodb"
)

// Config is the struct that contains the configuration for the service.
type Config struct {
	LogLevel string
	LogType  string
	Mongo    mongodb.Config
	Email    newsletter.EmailConfig
}

func main() {

	cfg := Config{
		LogLevel: getEnvWithDefault("LOG_LEVEL", ""),
		LogType:  getEnvWithDefault("LOG_TYPE", ""),
		Mongo: mongodb.Config{
			URI: getEnvWithDefault("NL_MONGO_URI", ""),
		},
		Email: newsletter.EmailConfig{
			Password: getEnvWithDefault("NL_EMAIL_PASSWORD", ""),
			Username: getEnvWithDefault("NL_EMAIL_USERNAME", ""),
		},
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	if err := setUpLog(cfg); err != nil {
		slog.Error("error setting up log", "error", err)
		signalCh <- syscall.SIGTERM
	}

	ctx := context.Background()

	client, err := mongodb.OpenDB(ctx, cfg.Mongo)
	if err != nil {
		slog.Error("error connecting to MongoDB", "error", err)
		signalCh <- syscall.SIGTERM
	}

	slog.Info("connected successfully to MongoDB instance")

	storage := mongodb.NewNLStorage(client, "newsletter")

	err = storage.SaveEngineer(ctx, mongodb.Engineer{
		Name:        "Paul Graham",
		URL:         "http://www.paulgraham.com/articles.html",
		Description: "Paul Graham is an English-born computer scientist, entrepreneur, venture capitalist, author, and essayist. He is best known for his work on Lisp, his former startup Viaweb (later renamed \"Yahoo! Store\"), co-founding the influential startup accelerator and seed capital firm Y Combinator, his blog, and Hacker News.",
	})
	if err != nil {
		slog.Error("error saving engineer", "error", err)
		signalCh <- syscall.SIGTERM
	}

	err = storage.SaveEngineer(ctx, mongodb.Engineer{
		Name:        "Joel Spolsky",
		URL:         "https://www.joelonsoftware.com/",
		Description: "Joel Spolsky is a software engineer and writer. He is the author of Joel on Software, a blog on software development, and the creator of the project management software Trello. He has previously worked as a programmer, software designer, and software consultant.",
	})
	if err != nil {
		slog.Error("error saving engineer", "error", err)
		signalCh <- syscall.SIGTERM
	}

	crawler := newsletter.NewCrawler(5, time.Duration(10)*time.Second, signalCh)

	go func() {
		crawler.Run(ctx, storage, newsletter.Fetch)
	}()

	mail := newsletter.NewMailClient(cfg.Email)

	go func() {
		for range time.Tick(time.Duration(50) * time.Second) {
			err := newsletter.EmailTrigger(ctx, storage, mail)
			if err != nil {
				slog.Error("error sending email", "error", err)
				signalCh <- syscall.SIGTERM
			}
		}
	}()

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

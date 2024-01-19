// Package main is the initial point for the service newsletter.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/perebaj/newsletter"
	"github.com/perebaj/newsletter/mongodb"
)

// Config is the struct that contains the configuration for the service.
type Config struct {
	LogLevel            string
	LogType             string
	LoopDurationMinutes time.Duration
	Mongo               mongodb.Config
}

func main() {

	cfg := Config{
		LogLevel: getEnvWithDefault("LOG_LEVEL", "INFO"),
		LogType:  getEnvWithDefault("LOG_TYPE", "json"),
		Mongo: mongodb.Config{
			URI: getEnvWithDefault("NL_MONGO_URI", ""),
		},
		LoopDurationMinutes: time.Duration(10) * time.Second,
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

	URLCh := make(chan string)
	fetchResultCh := make(chan string)

	var wg sync.WaitGroup
	wg.Add(5)

	for i := 0; i < 5; i++ {
		go newsletter.Worker(&wg, URLCh, fetchResultCh, newsletter.Fetch)
	}

	go func() {
		defer close(URLCh)
		for range time.Tick(cfg.LoopDurationMinutes) {
			slog.Info("fetching engineers")
			gotURLs, err := storage.DistinctEngineerURLs(ctx)
			if err != nil {
				slog.Error("error getting engineers", "error", err)
				signalCh <- syscall.SIGTERM
			}

			slog.Info("fetched engineers", "engineers", len(gotURLs))
			for _, url := range gotURLs {
				URLCh <- url.(string)
			}
		}
	}()

	go func() {
		wg.Wait()
		defer close(fetchResultCh)
	}()

	go func() {
		for v := range fetchResultCh {
			slog.Info("saving fetched sites response", "response", v[:10])
			err := storage.SaveSite(ctx, []mongodb.Site{
				{Content: v, ScrapeDatetime: time.Now().UTC()},
			})
			if err != nil {
				slog.Error("error saving site result", "error", err)
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

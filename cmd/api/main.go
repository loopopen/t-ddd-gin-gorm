package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/loopopen/pkg/confx"
	"github.com/loopopen/pkg/x"
	"github.com/loopopen/t-ddd-gin-gorm/cmd/api/config"
)

var commitSHA string
var f = flag.String("f", "etc/config.yaml", "config file")

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 15 * time.Second
	writeTimeout      = 15 * time.Second
	idleTimeout       = 60 * time.Second
)

func main() {
	flag.Parse()

	env := config.Env{
		Name:      os.Getenv("APP_ENV"),
		CommitSHA: commitSHA,
	}
	if env.Name == "" {
		env.Name = "prod"
	}

	var c config.Config
	confx.MustLoad(env.Name, *f, &c)

	if !env.IsProd() {
		cjson, _ := json.MarshalIndent(c, "", strings.Repeat(" ", 4))
		fmt.Printf("config of %s: %s\n", env.Name, string(cjson))
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	slogger := newLogger(c.Logger)

	app := x.Must(wireApp(&env, &c, c.Gap, c.Kafka, c.ORM, slogger))
	server := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", c.Bind, c.Port),
		Handler:           app.Engine,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slogger.Error("server stop listening with error", slog.Any("err", err))
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slogger.Info("shutting down gracefully, press Ctrl+C again to force")

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		slogger.Warn("forced shutdown by user, exiting immediately")
		os.Exit(1)
	}()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		slogger.Warn("server shutdown with error", slog.Any("err", err))
		os.Exit(1)
	}

	if err := app.ShutdownGap(); err != nil {
		slogger.Warn("gap shutdown with error", slog.Any("err", err))
		os.Exit(1)
	}

	slogger.Info("server shutdown successfully")
}

func newLogger(c config.Logger) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: c.LogLevel(),
	}
	var h slog.Handler
	if c.JSON {
		h = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		h = slog.NewTextHandler(os.Stdout, opts)
	}
	logger := slog.New(h)
	slog.SetDefault(logger)
	return logger
}

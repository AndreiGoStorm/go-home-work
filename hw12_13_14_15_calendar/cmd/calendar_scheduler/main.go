package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/config"
	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/logger"
	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/rabbit"
	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/scheduler"
	"github.com/AndreiGoStorm/go-home-work/hw12_13_14_15_calendar/internal/storage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.toml", "Path to configuration scheduler file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	conf := config.New(configFile)
	logg := logger.New(conf.Logger.Level)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	store := storage.New(conf)
	if err := store.Connect(ctx); err != nil {
		logg.Error("failed to store connect", err)
	}
	defer store.Close(ctx)

	rab := rabbit.New(conf, logg)
	if err := rab.Connect(); err != nil {
		logg.Error("failed to rabbit connect", err)
	}
	defer rab.Close()

	sch := scheduler.New(store, rab, logg)
	if err := sch.Run(ctx); err != nil {
		logg.Error("failed to run scheduler", err)
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

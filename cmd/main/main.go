package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type config struct {
	roomCode string
	delay    int
}

type app struct {
	server string
	logger *slog.Logger
	cfg    config
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	var cfg config

	flag.StringVar(&cfg.roomCode, "code", "", "room code")
	flag.IntVar(&cfg.delay, "delay", 2000, "delay between bot connections in ms")
	flag.Parse()

	if cfg.roomCode == "" {
		logger.Error("set room code with -code argument")
		os.Exit(1)
	}

	if cfg.delay < 200 {
		logger.Error("delay is too low", "delay", cfg.delay)
		os.Exit(1)
	}

	b := app{
		logger: logger,
		cfg:    cfg,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	b.logger.Info("running bots", "room code", b.cfg.roomCode, "delay", b.cfg.delay)
	if err := b.run(ctx); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

// if parsed 42[2,5,2]
// wait for 42[2,11,{"turnNum":0,"screen":3,"previous":null,"sentence":"Crazy singer hiding"}]
// and upload 42[2,6,{"t":0,"v":"gandon"}]
// 42[2,15,true]
// dont forget to send 2 to ping

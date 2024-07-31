package main

import (
	"context"
	"log/slog"
	"math/rand"
	"time"
)

func (b *bot) run(ctx context.Context) error {
	if err := b.findServer(); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			b.logger.Info("stopping bot")
			return nil
		case <-time.After(time.Millisecond * time.Duration(b.cfg.delay)):
			sid, err := b.getSid()
			if err != nil {
				return err
			}

			slog.Info("got session id", "sid", sid)

			if err := b.sendData(sid, rand.Intn(46)); err != nil {
				return err
			}
			slog.Info("joined the game", "sid", sid)
		}
	}
}

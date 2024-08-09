package main

import (
	"context"
	"time"

	"github.com/moonlags/garticBot/internal/entity"
)

func (app *app) run(ctx context.Context) error {
	if err := app.findServer(); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			app.logger.Info("stopping bots")
			return nil
		case <-time.After(time.Millisecond * time.Duration(app.cfg.delay)):
			// sid, err := b.getSid()
			// if err != nil {
			// 	return err
			// }

			bot := entity.NewBot(app.server, app.cfg.roomCode, app.cfg.name, app.logger)

			// go bot.Join()
			bot.Join()

			// slog.Info("got session id", "sid", sid)

			// if err := b.sendData(sid, rand.Intn(46)); err != nil {
			// 	return err
			// }
			// slog.Info("joined the game", "sid", sid)
		}
	}
}

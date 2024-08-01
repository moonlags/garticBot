package main

import (
	"github.com/levigross/grequests"
)

func (b *app) findServer() error {
	resp, err := grequests.Get("https://garticphone.com/api/server?code="+b.cfg.roomCode, &grequests.RequestOptions{UseCookieJar: true})
	if err != nil {
		return err
	}
	defer resp.Close()

	b.server = resp.String()
	b.logger.Info("found server", "server", b.server)

	return nil
}

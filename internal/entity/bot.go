package entity

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/levigross/grequests"
)

type Bot struct {
	Server   string
	Name     string
	Cookies  *http.CookieJar
	Logger   *slog.Logger
	RoomCode string
	ID       uuid.UUID
}

func NewBot(server string, roomCode string, name string, logger *slog.Logger) *Bot {
	return &Bot{
		ID:       uuid.New(),
		Server:   server,
		Logger:   logger,
		RoomCode: roomCode,
		Name:     name,
	}
}

func (b *Bot) Join() {
	b.Logger.Info("bot joining", "id", b.ID)

	sid, err := b.getSid()
	if err != nil {
		b.Logger.Error(err.Error(), "id", b.ID)
		return
	}

	if err := b.sendData(sid); err != nil {
		b.Logger.Error(err.Error(), "id", b.ID)
		return
	}
}

// TODO:
// get to server/socker.io/?EIO=3&transport=polling&sid=sid
// websocket to wss://server/socket.io/?EIO=3&transport=websocket&sid=sid
// get to server/socket.io/?EIO=3&transport=polling&sid=sid

func (b *Bot) sendData(sid string) error {
	payload := fmt.Sprintf("42[1,\"%s\",\"%s\",%d,\"en\",false,\"%s\",null,null]", b.ID, b.Name, rand.Intn(46), b.RoomCode)

	payload = fmt.Sprintf("%d:%s", len(payload), payload)
	r := strings.NewReader(payload)

	if _, err := grequests.Post(b.Server+"/socket.io/?EIO=3&transport=polling&sid="+sid, &grequests.RequestOptions{RequestBody: r, UseCookieJar: true}); err != nil {
		return err
	}

	return nil
}

func (b *Bot) getSid() (string, error) {
	resp, err := grequests.Get(b.Server+"/socket.io/?EIO=3&transport=polling", &grequests.RequestOptions{UseCookieJar: true})
	if err != nil {
		return "", err
	}
	defer resp.Close()

	data := resp.Bytes()
	data = data[bytes.IndexRune(data, '{'):]
	data = data[:bytes.LastIndexByte(data, '}')+1]

	var sid struct {
		Sid string `json:"sid"`
	}
	if err := json.Unmarshal(data, &sid); err != nil {
		return "", err
	}
	return sid.Sid, nil
}

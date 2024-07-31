package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/levigross/grequests"
)

// 88:42[1,"fed6ca58-9141-4c38-8d27-ebd3817aea70","dadas",6,"en",false,"0154940904",null,null]
func (b *bot) sendData(sid string, avatar int) error {
	payload := fmt.Sprintf("42[1,\"%s\",\"%s\",%d,\"en\",false,\"%s\",null,null]", uuid.NewString(), b.cfg.name, avatar, b.cfg.roomCode)
	payload = fmt.Sprintf("%d:%s", len(payload), payload)
	r := strings.NewReader(payload)

	if _, err := grequests.Post(b.server+"/socket.io/?EIO=3&transport=polling&sid="+sid, &grequests.RequestOptions{RequestBody: r, UseCookieJar: true}); err != nil {
		return err
	}

	return nil
}

func (b *bot) getSid() (string, error) {
	resp, err := grequests.Get(b.server+"/socket.io/?EIO=3&transport=polling", &grequests.RequestOptions{UseCookieJar: true})
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

func (b *bot) findServer() error {
	resp, err := grequests.Get("https://garticphone.com/api/server?code="+b.cfg.roomCode, &grequests.RequestOptions{UseCookieJar: true})
	if err != nil {
		return err
	}
	defer resp.Close()

	b.server = resp.String()
	b.logger.Info("found server", "server", b.server)

	return nil
}

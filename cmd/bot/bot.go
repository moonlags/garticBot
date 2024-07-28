package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// https://garticphone.com/en/?c=0154940904
func main() {
	code := flag.String("c", "", "room code")
	flag.Parse()

	if *code == "" {
		log.Fatal("Set room code with -c")
	}
	log.Println(*code)

	for {
		time.Sleep(1 * time.Second)
		client := http.DefaultClient
		// Санчос
		server, err := findServer(client, *code)
		if err != nil {
			log.Fatal("Failed to find server:", err)
		}
		log.Println(server)

		sid, err := getSid(client, server)
		if err != nil {
			log.Fatal("Failed to get session id:", err)
		}
		log.Println(sid)

		avatar := rand.Intn(46)
		name := base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(int(time.Now().Unix()))))

		if err := sendData(client, server, sid, *code, avatar, name); err != nil {
			log.Fatal("Failed to send user data:", err)
		}
	}
}

// 88:42[1,"fed6ca58-9141-4c38-8d27-ebd3817aea70","dadas",6,"en",false,"0154940904",null,null]
func sendData(client *http.Client, server string, sid string, code string, avatar int, name string) error {
	payload := fmt.Sprintf("42[1,\"%s\",\"%s\",%d,\"en\",false,\"%s\",null,null]", uuid.New(), name, avatar, code)
	payload = fmt.Sprintf("%d:%s", len(payload), payload)
	r := strings.NewReader(payload)

	resp, err := client.Post(server+"/socket.io/?EIO=3&transport=polling&sid="+sid, "text/plain;charset=UTF-8", r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func getSid(client *http.Client, server string) (string, error) {
	resp, err := client.Get(server + "/socket.io/?EIO=3&transport=polling")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

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

func findServer(client *http.Client, code string) (string, error) {
	resp, err := client.Get("https://garticphone.com/api/server?code=" + code)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	server, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(server), nil
}

// if parsed 42[2,5,2]
// wait for 42[2,11,{"turnNum":0,"screen":3,"previous":null,"sentence":"Crazy singer hiding"}]
// and upload 42[2,6,{"t":0,"v":"gandon"}]
// 42[2,15,true]
// dont forget to send 2 to ping

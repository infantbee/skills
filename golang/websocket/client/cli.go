package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var machID = flag.Int("id", 1111, "the machine id")

func main() {
	flag.Parse()

	url := fmt.Sprintf("ws://localhost:8080/ws?machid=%d", machID) //服务器地址
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			// ping
			err := ws.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Second * 2)
		}
	}()

	for {
		_, data, err := ws.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("receive: ", string(data))
	}
}

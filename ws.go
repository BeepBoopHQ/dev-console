package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	var err error
	wsConn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("Client connected: %s\n", wsConn.RemoteAddr().String())
	listen(wsConn)
}

func listen(wsConn *websocket.Conn) {
	defer wsConn.Close()

	for {
		_, b, err := wsConn.ReadMessage()
		if err != nil {
			if err != io.EOF {
				log.Println("NextReader:", err)
			}
			return
		}

		var f interface{}
		err = json.Unmarshal(b, &f)
		if err != nil {
			log.Printf("err: %s", err.Error())
		}

		msgMap := f.(map[string]interface{})
		if val, ok := msgMap["type"]; ok {
			if val == "auth" {
				wsConn.WriteJSON(NewAuthResultMessage(true, ""))

				// new connection, successfully auth'd, send all resources
				for _, botRes := range registry.BotResources {
					newAddRes := NewAddResourceMessage(idGen(), &Resource{
						ResourceCfg: botRes.ResourceCfg,
					})
					wsConn.WriteJSON(newAddRes)
				}
			}
		} else {
			fmt.Printf("Invalid message received: %s\n", string(b))
		}
	}
}

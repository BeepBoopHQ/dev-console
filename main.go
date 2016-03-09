package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var Port = "9000"
var addr = flag.String("addr", ":9000", "http service address")

// there is just one websocket connection at a time for this app
var wsConn *websocket.Conn

func main() {
	flag.Parse()

	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/api/resource", apiResourceHandler)
	http.HandleFunc("/css", cssHandler)
	http.HandleFunc("/jq", jQueryHandler)
	http.HandleFunc("/frm2js", frm2jsHandler)
	http.HandleFunc("/", viewHandler)

	log.Printf("Starting server at %#v", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

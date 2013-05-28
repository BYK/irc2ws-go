package main

import (
	"code.google.com/p/go.net/websocket"
	"flag"
	"log"
	"net/http"
)

var listenOn = flag.String("listen", ":1988", "Address to listen on")

func main() {
	flag.Parse()

	http.Handle("/", websocket.Handler(wsHandler))
	if err := http.ListenAndServe(*listenOn, nil); err != nil {
		log.Fatal("Error starting server:", err)
	}
}

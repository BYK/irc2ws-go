package main

import (
	"code.google.com/p/go.net/websocket"
	"flag"
	"log"
	"net/http"
)

func wsHandler(ws *websocket.Conn) {
	ircServerAddr := strings.TrimPrefix(ws.Request().URL.Path, "/")

	log.Println("Opening connection to ", ircServerAddr)
	ircConn, err := net.Dial("tcp", ircServerAddr)

	if err != nil {
		log.Println("Cannot open TCP connection to %s", ircServerAddr)
		ws.Close()
	} else {
		bridge := &WS2IRCBridge{ws: ws, irc: ircConn}
		bridge.run()
	}
}

var listenOn = flag.String("listen", ":1988", "Address to listen on")

func main() {
	flag.Parse()

	http.Handle("/", websocket.Handler(wsHandler))
	if err := http.ListenAndServe(*listenOn, nil); err != nil {
		log.Fatal("Error starting server:", err)
	}
}

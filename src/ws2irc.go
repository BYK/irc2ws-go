package main

import (
	"bufio"
	"code.google.com/p/go.net/websocket"
	"fmt"
	"net"
	"strings"
)

type WS2IRCBridge struct {
	ws  *websocket.Conn
	irc net.Conn
}

func (bridge *WS2IRCBridge) ws2irc() {
	writer := bufio.NewWriter(bridge.irc)

	for {
		var msg string
		wsErr := websocket.Message.Receive(bridge.ws, &msg)
		if wsErr != nil {
			fmt.Println("Error while reading from WebSocket: ", wsErr)
			break
		}

		_, ircErr := writer.WriteString(msg + "\r\n")
		if ircErr != nil {
			fmt.Println("Error while writing to IRC: ", ircErr)
			break
		}
		writer.Flush()
	}
	bridge.close()
}

func (bridge *WS2IRCBridge) irc2ws() {
	scanner := bufio.NewScanner(bridge.irc)

	for scanner.Scan() {
		msg := scanner.Text()
		err := websocket.Message.Send(bridge.ws, msg)
		if err != nil {
			fmt.Println("Error while writing to WebSocket: ", err)
			break
		}
	}
	bridge.close()
}

func (bridge *WS2IRCBridge) close() {
	bridge.ws.Close()
	bridge.irc.Close()
}

func (bridge *WS2IRCBridge) run() {
	go bridge.irc2ws()
	bridge.ws2irc()
}

func wsHandler(ws *websocket.Conn) {
	ircServerAddr := strings.TrimPrefix(ws.Request().URL.Path, "/")

	fmt.Println("Opening connection to ", ircServerAddr)
	ircConn, err := net.Dial("tcp", ircServerAddr)

	if err != nil {
		fmt.Println("Cannot open TCP connection to %s", ircServerAddr)
		ws.Close()
	} else {
		bridge := &WS2IRCBridge{ws: ws, irc: ircConn}
		bridge.run()
	}
}

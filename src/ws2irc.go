package main

import (
	"bufio"
	"code.google.com/p/go.net/websocket"
	"log"
	"net"
)

type WS2IRCBridge struct {
	ws  *websocket.Conn
	irc net.Conn
}

func (bridge *WS2IRCBridge) ws2irc() {
	defer bridge.close()
	writer := bufio.NewWriter(bridge.irc)

	for {
		var msg string
		wsErr := websocket.Message.Receive(bridge.ws, &msg)
		if wsErr != nil {
			log.Println("Error while reading from WebSocket: ", wsErr)
			break
		}

		_, ircErr := writer.WriteString(msg + "\r\n")
		if ircErr != nil {
			log.Println("Error while writing to IRC: ", ircErr)
			break
		}
		writer.Flush()
	}
}

func (bridge *WS2IRCBridge) irc2ws() {
	defer bridge.close()
	scanner := bufio.NewScanner(bridge.irc)

	for scanner.Scan() {
		msg := scanner.Text()
		err := websocket.Message.Send(bridge.ws, msg)
		if err != nil {
			log.Println("Error while writing to WebSocket: ", err)
			break
		}
	}
}

func (bridge *WS2IRCBridge) close() {
	bridge.ws.Close()
	bridge.irc.Close()
}

func (bridge *WS2IRCBridge) run() {
	go bridge.irc2ws()
	bridge.ws2irc()  // this function should be blocking
}

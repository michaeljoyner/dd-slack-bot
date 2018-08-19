package main

import (
	"log"
	"os"
	"strings"
	"sync/atomic"

	"github.com/michaeljoyner/dd-slack-bot/slack"
	"golang.org/x/net/websocket"
)

func main() {
	conn, err := slack.ConnectToRTM(os.Getenv("SLACK_DD_TOKEN"))

	if err != nil {
		log.Fatalf("failed to establish websocket: %v", err)
	}

	var counter uint64
	queue := make(chan slack.Message)
	heartbeat := make(chan slack.Ping)

	go func() {
		for {
			select {
			case m := <-queue:
				m.ID = atomic.AddUint64(&counter, 1)
				websocket.JSON.Send(conn.WS, m)
			case p := <-heartbeat:
				p.ID = atomic.AddUint64(&counter, 1)
				websocket.JSON.Send(conn.WS, p)
			}
		}
	}()

	conn.KeepAlive(heartbeat)

	for {
		mess, err := slack.GetMessage(conn.WS)
		if err != nil {
			log.Fatalf("could not get message: %v", err)
		}

		if mess.IsFor(conn.Self.ID) {
			cmd := strings.TrimPrefix(mess.Text, "<@"+conn.Self.ID+"> ")
			switch cmd {
			case "due":
				go handleDue(mess, queue)
			case "cost":
				go handleCost(mess, queue)
			default:
				go handleDefault(mess, queue)
			}
		}
	}
}

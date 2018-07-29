package main

import (
	"log"
	"os"
	"strings"

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

	go func() {
		for message := range queue {
			websocket.JSON.Send(conn.WS, message)
		}
	}()

	for {
		mess, err := slack.GetMessage(conn.WS)
		if err != nil {
			log.Fatalf("could not get message: %v", err)
		}

		if mess.IsFor(conn.Self.ID) {
			cmd := strings.TrimPrefix(mess.Text, "<@"+conn.Self.ID+"> ")
			switch cmd {
			case "due":
				go handleDue(mess, queue, counter)
			case "cost":
				go handleCost(mess, queue, counter)
			default:
				go handleDefault(mess, queue, counter)
			}
		}
	}
}

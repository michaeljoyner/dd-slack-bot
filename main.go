package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	"github.com/michaeljoyner/dd-slack-bot/dymantic"
	"golang.org/x/net/websocket"
)

type startResponse struct {
	Ok    bool         `json:"ok"`
	Error string       `json:"error"`
	URL   string       `json:"url"`
	Self  responseSelf `json:"self"`
}

type responseSelf struct {
	ID string `json:"id"`
}

//Message represents a slack message
type Message struct {
	ID      uint64 `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func (message *Message) isFor(id string) bool {
	return message.Type == "message" && strings.HasPrefix(message.Text, "<@"+id+">")
}

func getMessage(ws *websocket.Conn) (m Message, err error) {
	err = websocket.JSON.Receive(ws, &m)
	return
}

func main() {
	url := fmt.Sprintf("https://slack.com/api/rtm.start?token=%s", os.Getenv("SLACK_DD_TOKEN"))

	resp, err := http.Get(url)

	if err != nil {
		log.Fatalf("failed to get going: %v", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("failed to read body: %v", err)
	}

	var start startResponse

	jerr := json.Unmarshal(body, &start)

	if jerr != nil {
		log.Fatalf("failed to parse json: %v", err)
	}

	ws, err := websocket.Dial(start.URL, "", "https://api.slack.com/")

	if err != nil {
		log.Fatalf("failed to open websocket: %v", err)
	}

	var counter uint64

	for {
		mess, err := getMessage(ws)
		if err != nil {
			log.Fatalf("could not get message: %v", err)
		}

		if mess.isFor(start.Self.ID) {
			cmd := strings.TrimPrefix(mess.Text, "<@"+start.Self.ID+"> ")
			var m Message
			switch cmd {
			case "due":
				go func(m Message) {
					sites, err := dymantic.DueForHosting()
					if err != nil {
						fmt.Printf("could not fetch sites: %v", err)
					}
					for _, s := range sites {
						txt := fmt.Sprintf("*%v* is only paid until *%v*, they should be paying *%v*", s.Name, s.PaidUntil, s.HostingFee)
						m := Message{ID: atomic.AddUint64(&counter, 1), Type: "message", Channel: m.Channel, Text: txt}
						websocket.JSON.Send(ws, m)
					}
				}(mess)

			default:
				m = Message{ID: atomic.AddUint64(&counter, 1), Type: "message", Channel: mess.Channel, Text: "I don't know what you mean, but I like the way you say it."}
				websocket.JSON.Send(ws, m)
			}

		}
	}

}

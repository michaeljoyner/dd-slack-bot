package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

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

//Connection holds the websocket connection and the connected bots user
type Connection struct {
	WS   *websocket.Conn
	Self responseSelf
}

//IsFor determines if the message is intended for the given user ID
func (message *Message) IsFor(id string) bool {
	return message.Type == "message" && strings.HasPrefix(message.Text, "<@"+id+">")
}

//GetMessage reads a slack Message from the websocket
func GetMessage(ws *websocket.Conn) (m Message, err error) {
	err = websocket.JSON.Receive(ws, &m)
	return
}

//ConnectToRTM establishes a websocket connection to the Slack realtime api
func ConnectToRTM(token string) (Connection, error) {
	url := fmt.Sprintf("https://slack.com/api/rtm.start?token=%s", os.Getenv("SLACK_DD_TOKEN"))

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("failed to get going: %v", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Connection{}, err
	}

	var start startResponse

	if jerr := json.Unmarshal(body, &start); jerr != nil {
		return Connection{}, jerr
	}

	ws, err := websocket.Dial(start.URL, "", "https://api.slack.com/")

	if err != nil {
		return Connection{}, err
	}

	return Connection{WS: ws, Self: start.Self}, nil
}

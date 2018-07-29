package main

import (
	"fmt"
	"sync/atomic"

	"github.com/michaeljoyner/dd-slack-bot/dymantic"
	"github.com/michaeljoyner/dd-slack-bot/slack"
)

func handleDue(m slack.Message, q chan<- slack.Message, counter uint64) {
	sites, err := dymantic.DueForHosting()
	if err != nil {
		fmt.Printf("could not fetch sites: %v", err)
	}
	for _, s := range sites {
		txt := fmt.Sprintf("*%v* is only paid until *%v*, they should be paying *%v*", s.Name, s.PaidUntil, s.HostingFee)
		ms := slack.Message{ID: atomic.AddUint64(&counter, 1), Type: "message", Channel: m.Channel, Text: txt}
		q <- ms
	}
}

func handleCost(m slack.Message, q chan<- slack.Message, counter uint64) {
	summary, err := dymantic.TotalMonthlyHostingCost()
	if err != nil {
		fmt.Printf("could not fetch sites: %v", err)
	}
	txt := fmt.Sprintf("Current monthly hosting cost is *%v*", summary.PresentableCost)
	ms := slack.Message{ID: atomic.AddUint64(&counter, 1), Type: "message", Channel: m.Channel, Text: txt}
	q <- ms
}

func handleDefault(m slack.Message, q chan<- slack.Message, counter uint64) {
	ms := slack.Message{ID: atomic.AddUint64(&counter, 1), Type: "message", Channel: m.Channel, Text: "I don't know what you mean, but I like the way you say it."}
	q <- ms
}

package dymantic

import (
	"encoding/json"
	"os"
)

//Site represents a DD website we host
type Site struct {
	Name       string `json:"name"`
	HostingFee string `json:"presentable_hosting_fee"`
	Currency   string `json:"currency_name"`
	PaidUntil  string `json:"presentable_paid_until"`
}

//DueForHosting returns sites that require hosting payments in the next 31 days
func DueForHosting() ([]Site, error) {
	var sites []Site

	body, err := getWithAuthentication("/sites/due-payment", os.Getenv("SLACK_DD_TOKEN"))

	if err != nil {
		return nil, err
	}

	if jerr := json.Unmarshal(body, &sites); jerr != nil {
		return nil, jerr
	}

	return sites, nil

}

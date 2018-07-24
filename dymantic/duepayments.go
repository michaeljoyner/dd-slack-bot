package dymantic

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
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
	url := "http://secretadmin.dymanticdesign.com/admin-api/sites/due-payment"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authentication", "Bearer: "+os.Getenv("SLACK_DD_TOKEN"))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	var sites []Site
	jerr := json.Unmarshal(body, &sites)

	if jerr != nil {
		return nil, jerr
	}

	return sites, nil

}

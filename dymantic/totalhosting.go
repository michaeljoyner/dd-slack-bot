package dymantic

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

//HostingCostSummary summarises the cost of hosting for a month
type HostingCostSummary struct {
	TotalCost       int    `json:"total_cost"`
	Currency        string `json:"currency"`
	PresentableCost string `json:"presentable_cost"`
}

//TotalMonthlyHostingCost fetches a hosting cost summary
func TotalMonthlyHostingCost() (HostingCostSummary, error) {
	url := "http://secretadmin.dymanticdesign.com/admin-api/sites/hosting-cost"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return HostingCostSummary{}, err
	}

	req.Header.Add("Authentication", "Bearer: "+os.Getenv("SLACK_DD_TOKEN"))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return HostingCostSummary{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return HostingCostSummary{}, err
	}
	var summary HostingCostSummary

	jerr := json.Unmarshal(body, &summary)

	if jerr != nil {
		return HostingCostSummary{}, jerr
	}

	return summary, nil

}

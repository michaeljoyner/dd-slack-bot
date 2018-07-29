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

func getWithAuthentication(path, token string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, "http://secretadmin.dymanticdesign.com/admin-api"+path, nil)
	if err != nil {
		return []byte{}, err
	}

	req.Header.Add("Authentication", "Bearer: "+token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

//TotalMonthlyHostingCost fetches a hosting cost summary
func TotalMonthlyHostingCost() (HostingCostSummary, error) {
	var summary HostingCostSummary

	body, err := getWithAuthentication("/sites/hosting-cost", os.Getenv("SLACK_DD_TOKEN"))

	if err != nil {
		return HostingCostSummary{}, err
	}

	if jerr := json.Unmarshal(body, &summary); jerr != nil {
		return HostingCostSummary{}, jerr
	}

	return summary, nil

}

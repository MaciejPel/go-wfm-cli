package market

import (
	"encoding/json"
	"io"
	"math"
	"net/http"
)

// https://42bytes.notion.site/WFM-Api-v2-Documentation-5d987e4aa2f74b55a80db1a09932459d
const apiBase = "https://api.warframe.market/v2"

type ItemResponse struct {
	ApiVersion string   `json:"apiVersion"`
	Data       ItemData `json:"data"`
	Error      string   `json:"error"`
}

type ItemData struct {
	Sell  []OrderWithUser `json:"sell"`
	Buy   []OrderWithUser `json:"buy"`
	Error string          `json:"error,omitempty"`
}

type Order struct {
	Id         string `json:"id"`
	Type       string `json:"type"`
	Platinum   int32  `json:"platinum"`
	Quantity   int32  `json:"quantity"`
	PerTrade   int8   `json:"perTrade,omitempty"`
	Rank       *int8  `json:"rank,omitempty"`
	Charges    *int8  `json:"charges,omitempty"`
	Subtype    string `json:"subtype,omitempty"`
	AmberStars *int8  `json:"amberStars,omitempty"`
	CyanStars  *int8  `json:"cyanStars,omitempty"`
	Visible    bool   `json:"visible"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
	ItemId     string `json:"itemId"`
	Group      string `json:"group"`
}

type OrderWithUser struct {
	Order
}

type ItemPricing struct {
	Amt int
	Avg float64
	Min int
	Max int
}

func FetchItem(itemName string) (ItemPricing, error) {
	resp, err := http.Get(apiBase + "/orders/item/" + itemName + "/top")
	if err != nil {
		return ItemPricing{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return ItemPricing{}, err
	}

	var itemResponse ItemResponse
	err = json.Unmarshal(data, &itemResponse)
	if err != nil {
		return ItemPricing{}, err
	}

	options := []int{}
	for _, offer := range itemResponse.Data.Sell {
		options = append(options, int(offer.Platinum))
	}

	optionsLen := len(options)
	if optionsLen == 0 {
		return ItemPricing{}, err
	}

	return ItemPricing{
			Amt: optionsLen,
			Avg: avg(&options),
			Min: min(&options),
			Max: max(&options)},
		err
}

func min(arr *[]int) int {
	v := math.MaxInt
	for _, i := range *arr {
		if i < v {
			v = i
		}
	}
	return v
}

func max(arr *[]int) int {
	v := math.MinInt
	for _, i := range *arr {
		if i > v {
			v = i
		}
	}
	return v
}

func avg(arr *[]int) float64 {
	sum := 0
	for _, i := range *arr {
		sum += i
	}
	return float64(sum) / float64(len(*arr))
}

package market

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/MaciejPel/go-wfm-cli/pkg/constants"
	"github.com/MaciejPel/go-wfm-cli/pkg/utils"
)

// https://42bytes.notion.site/WFM-Api-v2-Documentation-5d987e4aa2f74b55a80db1a09932459d
const apiBase = "https://api.warframe.market/v2"

type ItemResponse struct {
	ApiVersion string   `json:"apiVersion"`
	Data       ItemData `json:"data"`
	Error      string   `json:"error"`
}

type ItemData struct {
	Sell  []Order `json:"sell"`
	Buy   []Order `json:"buy"`
	Error string  `json:"error,omitempty"`
}

type Order struct {
	Id       string `json:"id"`
	Type     string `json:"type"`
	Platinum int32  `json:"platinum"`
	Quantity int32  `json:"quantity"`
	PerTrade int8   `json:"perTrade,omitempty"`
	Visible  bool   `json:"visible"`
	ItemId   string `json:"itemId"`
}

type ItemJson struct {
	Id      string                   `json:"id"`
	Slug    string                   `json:"slug"`
	GameRef string                   `json:"gameRef"`
	Tags    []string                 `json:"tags,omitzero"`
	Ducats  int32                    `json:"ducats,omitempty"`
	I18N    map[string]*ItemI18NJson `json:"i18n,omitempty"`
}

type ItemI18NJson struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	WikiLink    string `json:"wikiLink,omitempty"`
	Icon        string `json:"icon"`
	Thumb       string `json:"thumb"`
	SubIcon     string `json:"subIcon,omitempty"`
}

type ItemPricing struct {
	Amt int
	Avg float64
	Min int
	Max int
}

type RelicItemsResponse struct {
	ApiVersion string     `json:"apiVersion"`
	Data       []ItemJson `json:"data"`
	Error      string     `json:"error"`
}

type ItemValue struct {
	Name   string
	Ducats int32
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
			Max: max(&options),
		},
		err
}

func FetchRelicItems(useCache bool) ([]ItemValue, error) {
	valid := []ItemValue{}

	if useCache {
		if _, err := os.Stat(constants.CacheFilePath); err == nil {
			content, err := os.ReadFile(constants.CacheFilePath)
			if err != nil {
				return valid, nil
			}
			for line := range strings.SplitSeq(string(content), "\n") {
				values := strings.Split(line, ";")
				ducats := 0
				if len(values) > 1 {
					parsed, err := strconv.Atoi(values[1])
					if err == nil {
						ducats = parsed
					}
				}
				valid = append(valid, ItemValue{Name: values[0], Ducats: int32(ducats)})
			}
			if len(valid) > 500 {
				return valid, nil
			}
		}
	}

	resp, err := http.Get(apiBase + "/items")
	if err != nil {
		return []ItemValue{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return []ItemValue{}, err
	}

	var relicItemsResponse RelicItemsResponse
	err = json.Unmarshal(data, &relicItemsResponse)
	if err != nil {
		return []ItemValue{}, err
	}

	for _, r := range relicItemsResponse.Data {
		name := r.I18N["en"].Name
		ducats := r.Ducats
		tags := r.Tags
		gameRef := r.GameRef

		set := slices.Contains(tags, "set")
		immortalMod := strings.Contains(gameRef, "/Immortal/Immortal") && !strings.Contains(gameRef, "Wildcard")

		if (ducats > 0 && !set) || immortalMod {
			valid = append(valid, ItemValue{Name: name, Ducats: ducats})
		}
	}

	valid = append(valid,
		ItemValue{Name: "Forma Blueprint", Ducats: 0},
		ItemValue{Name: "2 X Forma Blueprint", Ducats: 0},
		ItemValue{Name: "Riven Silver", Ducats: 0},
		ItemValue{Name: "1,200 X Kuva", Ducats: 0},
		ItemValue{Name: "Ayatan Amber Star", Ducats: 0},
		ItemValue{Name: "Exilus Weapon Adapter Blueprint", Ducats: 0},
	)
	utils.SaveStringToFile(constants.CacheFilePath, formatItems(valid))

	return valid, nil
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

func formatItems(items []ItemValue) string {
	var out strings.Builder

	for i, item := range items {
		out.WriteString(item.Name + ";" + fmt.Sprintf("%d", item.Ducats))
		if i < len(items)-1 {
			out.WriteString("\n")
		}
	}

	return out.String()
}

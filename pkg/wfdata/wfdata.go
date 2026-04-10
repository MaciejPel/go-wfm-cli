package wfdata

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strings"
	"unicode"

	"github.com/ulikunitz/xz/lzma"
)

const indexURL = "https://origin.warframe.com/PublicExport/index_en.txt.lzma"
const manifestURL = "https://content.warframe.com/PublicExport/Manifest/"

var cacheFilePath = os.TempDir() + "/wfdata-cache.txt"

var KubrowItems = map[string]string{
	"PrimeKubrowCollarABuckleComponent": "Kavasa Prime Buckle",
	"PrimeKubrowCollarABlueprint":       "Kavasa Prime Kubrow Collar Blueprint",
	"PrimeKubrowCollarABandComponent":   "Kavasa Prime Band",
}
var ArchwingItems = map[string]string{
	"PrimeArchwingBlueprint":        "Odonata Prime Blueprint",
	"PrimeArchwingSystemsBlueprint": "Odonata Prime Systems Blueprint",
	"PrimeArchwingChassisBlueprint": "Odonata Prime Chassis Blueprint",
	"PrimeArchwingWingsBlueprint":   "Odonata Prime Wings Blueprint",
}
var ImmortalMods = map[string]string{
	"ImmortalOneMod":   "Lohk",
	"ImmortalTwoMod":   "Xata",
	"ImmortalThreeMod": "Jahu",
	"ImmortalFourMod":  "Vome",
	"ImmortalFiveMod":  "Ris",
	"ImmortalSixMod":   "Fass",
	"ImmortalSevenMod": "Netra",
	"ImmortalEightMod": "Khra",
}

type Resources struct {
	Resources []Resource `json:"ExportResources"`
}
type Resource struct {
	UniqueName        string `json:"uniqueName"`
	Name              string `json:"name"`
	PrimeSellingPrice *int   `json:"primeSellingPrice,omitempty"`
}
type Relics struct {
	Relics []Relic `json:"ExportRelicArcane"`
}
type Relic struct {
	UniqueName string   `json:"uniqueName"`
	Name       string   `json:"name"`
	Rewards    []Reward `json:"relicRewards"`
}
type Reward struct {
	Name string `json:"rewardName"`
}

func GetWarframeData(useCache bool) ([]string, error) {
	valid := []string{}

	if useCache {
		if _, err := os.Stat(cacheFilePath); err == nil {
			content, err := os.ReadFile(cacheFilePath)
			if err != nil {
				return valid, nil
			}
			valid = strings.Split(string(content), "\n")
			if len(valid) > 500 {
				return valid, nil
			}
		}
	}

	resp, err := http.Get(indexURL)
	if err != nil {
		return valid, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return valid, err
	}

	reader, err := lzma.NewReader(bytes.NewReader(data))
	if err != nil {
		return valid, err
	}

	decoded, err := io.ReadAll(reader)
	if err != nil {
		return valid, err
	}

	lines := strings.Split(strings.TrimSpace(string(decoded)), "\r\n")

	resources := make(map[string]string)
	rewards := []string{}

	for _, line := range lines {
		matched, err := regexp.Match(`Export(Resources|RelicArcane)_en.json.*`, []byte(line))
		if err != nil {
			return valid, err
		}
		if !matched {
			continue
		}

		resp, err := http.Get(manifestURL + line)
		if err != nil {
			return valid, err
		}
		defer resp.Body.Close()

		content, err := io.ReadAll(resp.Body)
		if err != nil {
			return valid, err
		}

		if strings.HasPrefix(line, "ExportResources_en.json") {
			var resourcesWrapper Resources
			err = json.Unmarshal(content, &resourcesWrapper)
			if err != nil {
				return valid, err
			}

			for _, r := range resourcesWrapper.Resources {
				if r.PrimeSellingPrice != nil {
					resources[r.UniqueName] = r.Name
				}
			}
		}

		if strings.HasPrefix(line, "ExportRelicArcane_en.json") {
			var relicsWrapper Relics
			err = json.Unmarshal(content, &relicsWrapper)
			if err != nil {
				return valid, err
			}

			for _, relics := range relicsWrapper.Relics {
				for _, reward := range relics.Rewards {
					if !slices.Contains(rewards, reward.Name) && (strings.Contains(reward.Name, "Prime") || strings.Contains(reward.Name, "Immortal")) {
						rewards = append(rewards, reward.Name)
					}
				}
			}
		}

	}

	for _, reward := range rewards {
		rewardSplit := strings.Split(reward, "/")
		lastItem := rewardSplit[len(rewardSplit)-1]

		if slices.Contains(rewardSplit, "WarframeRecipes") {
			if strings.HasSuffix(reward, "PrimeBlueprint") {
				valid = append(valid, separateByUppercase(lastItem))
			} else {
				valid = append(valid, strings.ReplaceAll(separateByUppercase(lastItem), "Helmet", "Neuroptics"))
			}
		} else if slices.Contains(rewardSplit, "SentinelRecipes") {
			re := regexp.MustCompile(`Prime|Sentinel|Blueprint`)
			sentinel := re.ReplaceAllString(lastItem, "")
			valid = append(valid, sentinel+" Prime Blueprint")
		} else if slices.Contains(rewardSplit, "Kubrow") {
			valid = append(valid, KubrowItems[lastItem])
		} else if slices.Contains(rewardSplit, "ArchwingRecipes") {
			valid = append(valid, ArchwingItems[lastItem])
		} else if slices.Contains(rewardSplit, "Immortal") {
			valid = append(valid, ImmortalMods[lastItem])
		} else if slices.Contains(rewardSplit, "WeaponParts") {
			weaponKey := strings.ReplaceAll(reward, "/StoreItems", "")
			valid = append(valid, resources[weaponKey])
		} else if slices.Contains(rewardSplit, "Weapons") {
			valid = append(valid, separateByUppercase(lastItem))
		}
	}

	saveStringToFile(cacheFilePath, strings.Join(valid, "\n"))

	return valid, nil
}

func saveStringToFile(path string, data string) {
	file, err := os.Create(path)
	if err != nil {
		return
	}
	defer file.Close()
	file.WriteString(data)
}

func separateByUppercase(input string) string {
	var out strings.Builder
	for i, v := range input {
		if unicode.IsUpper(v) && i > 0 {
			out.WriteString(" " + string(v))
		} else {
			out.WriteString(string(v))
		}
	}
	return out.String()
}

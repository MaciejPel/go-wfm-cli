package warframe

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/MaciejPel/go-wfm-cli/pkg/constants"
	"github.com/MaciejPel/go-wfm-cli/pkg/text"
	"github.com/MaciejPel/go-wfm-cli/pkg/utils"
	"github.com/ulikunitz/xz/lzma"
)

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

func GetData(useCache bool) ([]string, error) {
	valid := []string{}

	if useCache {
		if _, err := os.Stat(constants.CacheFilePath); err == nil {
			content, err := os.ReadFile(constants.CacheFilePath)
			if err != nil {
				return valid, nil
			}
			valid = strings.Split(string(content), "\n")
			if len(valid) > 500 {
				return valid, nil
			}
		}
	}

	resp, err := http.Get(constants.WfIndexURL)
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

		resp, err := http.Get(constants.WfManifestURL + line)
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
				valid = append(valid, text.SeparateByUppercase(lastItem))
			} else {
				valid = append(valid, strings.ReplaceAll(text.SeparateByUppercase(lastItem), "Helmet", "Neuroptics"))
			}
		} else if slices.Contains(rewardSplit, "SentinelRecipes") {
			re := regexp.MustCompile(`Prime|Sentinel|Blueprint`)
			sentinel := re.ReplaceAllString(lastItem, "")
			valid = append(valid, sentinel+" Prime Blueprint")
		} else if slices.Contains(rewardSplit, "Kubrow") {
			valid = append(valid, constants.KubrowItems[lastItem])
		} else if slices.Contains(rewardSplit, "ArchwingRecipes") {
			valid = append(valid, constants.ArchwingItems[lastItem])
		} else if slices.Contains(rewardSplit, "Immortal") {
			valid = append(valid, constants.ImmortalMods[lastItem])
		} else if slices.Contains(rewardSplit, "WeaponParts") {
			weaponKey := strings.ReplaceAll(reward, "/StoreItems", "")
			valid = append(valid, resources[weaponKey])
		} else if slices.Contains(rewardSplit, "Weapons") {
			valid = append(valid, text.SeparateByUppercase(lastItem))
		}
	}

	valid = append(valid, "Forma Blueprint")
	utils.SaveStringToFile(constants.CacheFilePath, strings.Join(valid, "\n"))

	return valid, nil
}

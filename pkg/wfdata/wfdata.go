package wfdata

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/ulikunitz/xz/lzma"
)

const indexURL = "https://origin.warframe.com/PublicExport/index_en.txt.lzma"
const manifestURL = "https://content.warframe.com/PublicExport/Manifest/"

// type Warframes struct {
// 	ExportWarframes []Warframe `json:"ExportWarframes"`
// }

// type Warframe struct {
// 	Name string `json:"name"`
// }

// type Weapons struct {
// 	ExportWeapons []Weapon `json:"ExportWeapons"`
// }

// type Weapon struct {
// 	Name string `json:"name"`
// }

type Resources struct {
	ExportResources []Resource `json:"ExportResources"`
}

type Resource struct {
	Name              string `json:"name"`
	PrimeSellingPrice *int   `json:"primeSellingPrice,omitempty"`
}

func Fetch() error {
	resp, err := http.Get(indexURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	reader, err := lzma.NewReader(bytes.NewReader(data))
	if err != nil {
		return err
	}

	decoded, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	var lines []string
	scanner := bufio.NewScanner(bytes.NewReader(decoded))
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}

	for _, line := range lines {
		url := manifestURL + line

		if /*strings.HasPrefix(line, "ExportWarframes_en.json") || strings.HasPrefix(line, "ExportWeapons_en.json") ||*/ strings.HasPrefix(line, "ExportResources_en.json") {

			resp, err := http.Get(url)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			content, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			// if strings.HasPrefix(line, "ExportWarframes_en.json") {
			// 	var wf Warframes
			// 	err = json.Unmarshal(content, &wf)
			// 	if err != nil {
			// 		return err
			// 	}

			// 	for _, x := range wf.ExportWarframes {
			// 		if strings.Contains(x.Name, "Prime") {
			// 			fmt.Println(x.Name)
			// 		}
			// 	}
			// }

			// if strings.HasPrefix(line, "ExportWeapons_en.json") {
			// 	var wf Weapons
			// 	err = json.Unmarshal(content, &wf)
			// 	if err != nil {
			// 		return err
			// 	}

			// 	for _, x := range wf.ExportWeapons {
			// 		if strings.Contains(x.Name, "Prime") {
			// 			fmt.Println(x.Name)
			// 		}
			// 	}
			// }

			if strings.HasPrefix(line, "ExportResources_en.json") {
				var wf Resources
				err = json.Unmarshal(content, &wf)
				if err != nil {
					return err
				}

				f := []string{}
				for _, resource := range wf.ExportResources {
					if strings.Contains(resource.Name, "Prime") && resource.PrimeSellingPrice != nil {
						if strings.Contains(resource.Name, "Chassis") {
							f = append(f, strings.ReplaceAll(resource.Name, "Chassis", "Blueprint"))
						}
						f = append(f, resource.Name)
					}
				}

				d1 := strings.Join(f, "\n")
				file, err := os.Create("data.txt")
				if err != nil {
					panic(err)
				}
				file.WriteString(d1)
				defer file.Close()

			}

		}

	}

	return nil
}

package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"

	c "github.com/MaciejPel/go-wfm-cli/pkg/constants"
	"github.com/MaciejPel/go-wfm-cli/pkg/imageutil"
	"github.com/MaciejPel/go-wfm-cli/pkg/market"
	"github.com/MaciejPel/go-wfm-cli/pkg/text"
	"github.com/MaciejPel/go-wfm-cli/pkg/utils"
	"github.com/eiannone/keyboard"
	"github.com/kbinani/screenshot"
	"github.com/tiagomelo/go-ocr/ocr"
)

func main() {
	validEntries, err := market.FetchRelicItems(true)
	if err != nil {
		log.Fatal(err)
	}

	err = keyboard.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer keyboard.Close()

	t, err := ocr.New(ocr.TesseractPath(c.TesseractPath))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	utils.ClearScreen()

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}

		if key == keyboard.KeyEsc || key == keyboard.KeyCtrlC || char == 'q' {
			fmt.Println("Bye!")
			break
		}

		if char == 'u' {
			validEntries, err = market.FetchRelicItems(false)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Item data updated")
		}

		if char == 'c' {
			utils.ClearScreen()
		}

		if char == 's' {
			utils.ClearScreen()

			bounds := screenshot.GetDisplayBounds(0)
			screenshotImg, err := screenshot.CaptureRect(bounds)
			err = imageutil.SavePNG(c.TmpImgPath, screenshotImg)
			if err != nil {
				panic(err)
			}

			img, err := imageutil.Load(c.TmpImgPath)
			if err != nil {
				panic(err)
			}

			rewardIndicators := []int{}
			for rewardAmt := 4; rewardAmt >= 1; rewardAmt-- {
				for rewardNo := range rewardAmt {
					centerX := c.BoxCropXStart[rewardAmt] + rewardNo*c.Box.Width + rewardNo*c.Box.Gap + c.Box.Width/2
					b := image.Rect(centerX-5, c.Box.IndicatorYStart, centerX+5, c.Box.IndicatorYEnd).Bounds()
					typeCount := []int{0, 0, 0}
					for x := b.Min.X; x < b.Max.X; x++ {
						for y := b.Min.Y; y < b.Max.Y; y++ {
							r, g, b, _ := img.At(x, y).RGBA()
							p := []int{int(r >> 8), int(g >> 8), int(b >> 8)}
							for index, color := range c.RewardIndicatorColors {
								if p[0] == color[0] && p[1] == color[1] && p[2] == color[2] {
									typeCount[index]++
								}
							}
						}
					}
					if rewardNo == 0 && (typeCount[0] < 10 && typeCount[1] < 10 && typeCount[2] < 10) {
						break
					}
					for indicator, count := range typeCount {
						if count >= 10 {
							rewardIndicators = append(rewardIndicators, indicator)
							break
						}
					}
				}
				if len(rewardIndicators) == rewardAmt {
					break
				}
			}

			if len(rewardIndicators) == 0 {
				fmt.Println("no rewards detected")
				continue
			}

			fmt.Printf("%-40s - %3s %4s %6s\n", "Item", "duc", "min", "avg")
			for i, indicator := range rewardIndicators {
				cropImgPath := os.TempDir() + "/wf-data-tmp-img-crop-" + strconv.Itoa(i) + ".png"
				startX := c.BoxCropXStart[len(rewardIndicators)] + i*c.Box.Gap + i*c.Box.Width
				rect := image.Rect(startX+c.Box.Border, c.Box.LabelYStart, startX+c.Box.Width-c.Box.Border, c.Box.LabelYEnd)

				cropped := imageutil.Crop(img, rect)
				gray := imageutil.ApplyGrayscale(cropped)
				otsu := imageutil.ApplyOtsuThreshold(gray)
				binaryInv := imageutil.ApplyThresholdBinaryInv(gray, otsu)

				err = imageutil.SavePNG(cropImgPath, binaryInv)
				if err != nil {
					panic(err)
				}

				extractedText, err := t.TextFromImageFile(cropImgPath)
				if err != nil {
					log.Fatal(err)
					os.Exit(1)
				}

				re := regexp.MustCompile(`[^0-9a-zA-Z&, ]+`)
				result := re.ReplaceAllString(strings.ReplaceAll(extractedText, "\r\n", ""), "")

				possibleEntries := []market.ItemValue{}
				for _, entry := range validEntries {
					if slices.Contains(c.DucatsToIndicators[int(entry.Ducats)], indicator) {
						possibleEntries = append(possibleEntries, entry)
					}
				}

				bestMatch := text.BestMatch(result, validEntries)
				if bestMatch.Ducats == 0 {
					fmt.Println(bestMatch.Name)
					continue
				}
				out, err := market.FetchItem(text.ItemNameToSlug(bestMatch.Name))
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("%-40s - %3d %4d %6.2f\n", bestMatch.Name, bestMatch.Ducats, out.Min, out.Avg)
			}
		}
	}

	os.Exit(0)
}

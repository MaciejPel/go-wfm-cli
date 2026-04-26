package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/MaciejPel/go-wfm-cli/pkg/constants"
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

	t, err := ocr.New(ocr.TesseractPath(constants.TesseractPath))
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
			err = imageutil.SavePNG(constants.TmpImgPath, screenshotImg)
			if err != nil {
				panic(err)
			}

			img, err := imageutil.Load(constants.TmpImgPath)
			if err != nil {
				panic(err)
			}

			fmt.Printf("%-40s - %4s %6s\n", "Item", "min", "avg")
			playerCount := 4
			for i := range playerCount {
				cropImgPath := os.TempDir() + "/wf-data-tmp-img-crop-" + strconv.Itoa(i) + ".png"
				boxTextXStart := constants.BoxCropXStart[playerCount] + i*constants.BoxGapWidth + i*constants.BoxWidth
				rect := image.Rect(boxTextXStart, constants.BoxTextYStart, boxTextXStart+constants.BoxWidth, constants.BoxTextYEnd)

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

				re := regexp.MustCompile(`[^a-zA-Z]+`)
				result := re.ReplaceAllString(strings.ReplaceAll(extractedText, "\r\n", ""), "")
				bestMatch := text.BestMatch(result, validEntries)
				if bestMatch == "Forma Blueprint" {
					fmt.Println(bestMatch)
					continue
				}
				out, err := market.FetchItem(text.ItemNameToSlug(bestMatch))
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("%-40s - %4d %6.2f\n", bestMatch, out.Min, out.Avg)
			}
		}
	}

	os.Exit(0)
}

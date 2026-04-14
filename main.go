package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/MaciejPel/go-wfm-cli/pkg/constants"
	"github.com/MaciejPel/go-wfm-cli/pkg/text"
	"github.com/MaciejPel/go-wfm-cli/pkg/warframe"
	"github.com/eiannone/keyboard"
	"github.com/kbinani/screenshot"
	"github.com/tiagomelo/go-ocr/ocr"
	"gocv.io/x/gocv"
)

func main() {
	validEntries, fetchErr := warframe.GetData(true)
	if fetchErr != nil {
		log.Fatal(fetchErr)
	}

	err := keyboard.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer keyboard.Close()

	fmt.Println("Listening...")
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}

		if key == keyboard.KeyEsc || key == keyboard.KeyCtrlC || char == 'q' {
			break
		}

		if char == 'u' {
			validEntries, err = warframe.GetData(false)
			if err != nil {
				log.Fatal(err)
			}
		}

		if char == 's' {
			bounds := screenshot.GetDisplayBounds(0)
			imga, err := screenshot.CaptureRect(bounds)
			if err != nil {
				panic(err)
			}
			file, _ := os.Create(constants.TmpImgPath)
			defer file.Close()
			png.Encode(file, imga)

			img := gocv.IMRead(constants.TmpImgPath, gocv.IMReadColor)
			// pixel := img.GetVecbAt(115, 350)
			if img.Empty() {
				panic("Cannot read image")
			}
			defer img.Close()

			for i, e := range constants.CropRegions[4] {
				cropImgPath := os.TempDir() + "/wf-data-tmp-img-crop-" + strconv.Itoa(i) + ".jpg"
				rect := image.Rect(e[0], e[1], e[2], e[3])
				cropped := img.Region(rect)
				defer cropped.Close()
				gray := gocv.NewMat()
				defer gray.Close()
				gocv.CvtColor(cropped, &gray, gocv.ColorBGRToGray)
				bw := gocv.NewMat()
				defer bw.Close()
				gocv.Threshold(gray, &bw, 0, 255, gocv.ThresholdBinaryInv|gocv.ThresholdOtsu)
				gocv.IMWrite(cropImgPath, bw)

				t, err := ocr.New(ocr.TesseractPath(constants.TesseractPath))
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				extractedText, err := t.TextFromImageFile(cropImgPath)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				re := regexp.MustCompile(`[^a-zA-Z]+`)
				result := re.ReplaceAllString(strings.ReplaceAll(extractedText, "\r\n", ""), "")
				fmt.Println(text.BestMatch(result, validEntries))
			}
		}
	}

}

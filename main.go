package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/tiagomelo/go-ocr/ocr"
	"gocv.io/x/gocv"
)

const tesseractPath = "C:\\Program Files\\Tesseract-OCR\\tesseract.exe"
const imagePath = "C:\\Users\\admin\\Desktop\\Untitled0.png"

func main() {
	err := keyboard.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer keyboard.Close()

	img := gocv.IMRead(imagePath, gocv.IMReadColor)
	if img.Empty() {
		panic("Cannot read image")
	}
	defer img.Close()
	rect := image.Rect(635, 540, 1920, 612)
	cropped := img.Region(rect)
	defer cropped.Close()
	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(cropped, &gray, gocv.ColorBGRToGray)
	bw := gocv.NewMat()
	defer bw.Close()
	gocv.Threshold(gray, &bw, 0, 255, gocv.ThresholdBinary|gocv.ThresholdOtsu)
	gocv.IMWrite("output.jpg", bw)

	fmt.Println("Listening...")
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}

		if key == keyboard.KeyEsc || key == keyboard.KeyCtrlC || char == 'q' {
			break
		}

		if char == 's' {
			t, err := ocr.New(ocr.TesseractPath(tesseractPath))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			extractedText, err := t.TextFromImageFile("./output.jpg")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			split := strings.Split(extractedText, "\n")
			for idx, val := range split {
				fmt.Println(idx, val)
			}
		}

	}

}

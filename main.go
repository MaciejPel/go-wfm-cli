package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/tiagomelo/go-ocr/ocr"
)

const tesseractPath = "/usr/bin/tesseract"
const imagePath = "/home/mp/Desktop/Untitled0.png"

func main() {
	defer keyboard.Close()
	err := keyboard.Open()
	if err != nil {
		log.Fatal(err)
	}

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

			extractedText, err := t.TextFromImageFile(imagePath)
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

package imageutil

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

func Load(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func Crop(img image.Image, rect image.Rectangle) image.Image {
	cropped := image.NewRGBA(rect)
	draw.Draw(cropped, rect, img, rect.Min, draw.Src)
	return cropped
}

func ApplyGrayscale(img image.Image) *image.Gray {
	bounds := img.Bounds()
	gray := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			grayColor := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			gray.Set(x, y, grayColor)
		}
	}
	return gray
}

func ApplyOtsuThreshold(img *image.Gray) uint8 {
	var hist [256]int
	total := 0

	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			val := img.GrayAt(x, y).Y
			hist[val]++
			total++
		}
	}

	var sum float64
	for i := range 256 {
		sum += float64(i * hist[i])
	}

	var sumB float64
	var wB, wF int
	var maxVar float64
	var threshold uint8

	for t := range 256 {
		wB += hist[t]
		if wB == 0 {
			continue
		}

		wF = total - wB
		if wF == 0 {
			break
		}

		sumB += float64(t * hist[t])
		mB := sumB / float64(wB)
		mF := (sum - sumB) / float64(wF)
		betweenVar := float64(wB*wF) * (mB - mF) * (mB - mF)

		if betweenVar > maxVar {
			maxVar = betweenVar
			threshold = uint8(t)
		}
	}

	return threshold
}

func ApplyThresholdBinaryInv(img *image.Gray, thresh uint8) *image.Gray {
	bounds := img.Bounds()
	out := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			val := img.GrayAt(x, y).Y
			if val > thresh {
				out.SetGray(x, y, color.Gray{Y: 0})
			} else {
				out.SetGray(x, y, color.Gray{Y: 255})
			}
		}
	}
	return out
}

func SavePNG(path string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

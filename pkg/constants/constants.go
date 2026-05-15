package constants

import "os"

const TesseractPath = "C:\\Program Files\\Tesseract-OCR\\tesseract.exe"

var CacheFilePath = os.TempDir() + "/wfdata-cache.txt"
var TmpImgPath = os.TempDir() + "/wf-data-tmp-img.png"

var Box = struct {
	Width           int
	Gap             int
	Border          int
	LabelYStart     int
	LabelYEnd       int
	IndicatorYStart int
	IndicatorYEnd   int
}{
	Width:           320,
	Gap:             3,
	Border:          4,
	LabelYStart:     551,
	LabelYEnd:       612,
	IndicatorYStart: 630,
	IndicatorYEnd:   650,
}

var BoxCropXStart = map[int]int{
	1: 1120,
	2: 958,
	3: 797,
	4: 635,
}

var RewardIndicatorColors = [][]int{
	{79, 47, 33},  //bronze
	{90, 89, 90},  //silver
	{101, 85, 30}, //gold
}

var DucatsToIndicators = map[int][]int{
	0:   {0, 1, 2},
	15:  {0},
	25:  {0, 1},
	45:  {1},
	65:  {1, 2},
	100: {2},
}

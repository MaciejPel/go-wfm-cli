package constants

import "os"

const TesseractPath = "C:\\Program Files\\Tesseract-OCR\\tesseract.exe"
const WfIndexURL = "https://origin.warframe.com/PublicExport/index_en.txt.lzma"
const WfManifestURL = "https://content.warframe.com/PublicExport/Manifest/"

var CacheFilePath = os.TempDir() + "/wfdata-cache.txt"
var TmpImgPath = os.TempDir() + "/wf-data-tmp-img.png"

const BoxWidth = 313
const BoxGapWidth = 10
const BoxTextYStart = 551
const BoxTextYEnd = 612

var BoxCropXStart = map[int]int{
	1: 1123,
	2: 961,
	3: 800,
	4: 638,
}

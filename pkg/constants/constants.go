package constants

import "os"

const TesseractPath = "C:\\Program Files\\Tesseract-OCR\\tesseract.exe"
const WfIndexURL = "https://origin.warframe.com/PublicExport/index_en.txt.lzma"
const WfManifestURL = "https://content.warframe.com/PublicExport/Manifest/"

var CacheFilePath = os.TempDir() + "/wfdata-cache.txt"
var TmpImgPath = os.TempDir() + "/wf-data-tmp-img.png"

var KubrowItems = map[string]string{
	"PrimeKubrowCollarABuckleComponent": "Kavasa Prime Buckle",
	"PrimeKubrowCollarABlueprint":       "Kavasa Prime Kubrow Collar Blueprint",
	"PrimeKubrowCollarABandComponent":   "Kavasa Prime Band",
}
var ArchwingItems = map[string]string{
	"PrimeArchwingBlueprint":        "Odonata Prime Blueprint",
	"PrimeArchwingSystemsBlueprint": "Odonata Prime Systems Blueprint",
	"PrimeArchwingChassisBlueprint": "Odonata Prime Chassis Blueprint",
	"PrimeArchwingWingsBlueprint":   "Odonata Prime Wings Blueprint",
}
var ImmortalMods = map[string]string{
	"ImmortalOneMod":   "Lohk",
	"ImmortalTwoMod":   "Xata",
	"ImmortalThreeMod": "Jahu",
	"ImmortalFourMod":  "Vome",
	"ImmortalFiveMod":  "Ris",
	"ImmortalSixMod":   "Fass",
	"ImmortalSevenMod": "Netra",
	"ImmortalEightMod": "Khra",
}

var CropRegions = map[int][][]int{
	4: [][]int{
		{638, 551, 951, 612},
		{961, 551, 1274, 612},
		{1284, 551, 1597, 612},
		{1607, 551, 1921, 612},
	},
}

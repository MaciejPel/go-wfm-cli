package utils

import (
	"os"
)

func SaveStringToFile(path string, data string) {
	file, err := os.Create(path)
	if err != nil {
		return
	}
	defer file.Close()
	file.WriteString(data)
}

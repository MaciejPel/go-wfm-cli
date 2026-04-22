package utils

import (
	"fmt"
	"os"
	"os/exec"
)

func SaveStringToFile(path string, data string) {
	file, err := os.Create(path)
	if err != nil {
		return
	}
	defer file.Close()
	file.WriteString(data)
}

func ClearScreen() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
	fmt.Println("Listening...")
}

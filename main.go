package main

import (
	"fmt"
	"os"
	"strings"
)

type config struct {
	Bytes int
}

func main() {

	cfg := config{
		Bytes: 8,
	}

	file, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println("Could not open file:", err)
		return
	}

	currentLine := ""
	readFile := make([]byte, cfg.Bytes)

	for {
		n, err := file.Read(readFile)
		if err != nil {
			return
		}

		splitBytes, splitLength := splitNewLine(string(readFile[:n]))
		currentLine += splitBytes[0]
		if splitLength > 1 {
			fmt.Printf("read: %s\n", currentLine)
			currentLine = splitBytes[1]
		}

	}

}

func splitNewLine(part string) ([]string, int) {
	result := strings.Split(part, "\n")
	return result, len(result)
}

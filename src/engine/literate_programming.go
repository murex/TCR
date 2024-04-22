package engine

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func makeLiterate(inputFile string) {
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("Error opening source file:", err)
		return
	}
	defer file.Close()

	outputFile := inputFile + ".md"
	outFile, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outFile.Close()

	scanner := bufio.NewScanner(file)
	inCommentBlock := false
	codeBlock := ""
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "/*") {
			inCommentBlock = true
			line = strings.TrimPrefix(line, "/*")
		}
		if strings.HasSuffix(line, "*/") {
			inCommentBlock = false
			line = strings.TrimSuffix(line, "*/")
		}
		if strings.HasPrefix(line, "//") || inCommentBlock {
			if codeBlock != "" {
				_, err := outFile.WriteString("```go\n" + codeBlock + "```\n")
				if err != nil {
					fmt.Println("Error writing to output file:", err)
					return
				}
				codeBlock = ""
			}
			line = strings.TrimPrefix(line, "//")
			_, err := outFile.WriteString(line + "\n")
			if err != nil {
				fmt.Println("Error writing to output file:", err)
				return
			}
		} else {
			codeBlock += line + "\n"
		}
	}
	if codeBlock != "" {
		_, err := outFile.WriteString("```go\n" + codeBlock + "```\n")
		if err != nil {
			fmt.Println("Error writing to output file:", err)
			return
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning source file:", err)
	}
}

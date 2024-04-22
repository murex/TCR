package engine

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type TourGuide struct {
	TourName    string
	StepTitle   string
	CodeSnippet string
	LineNumber  int
}

func extractTourGuides(sourceFile, outputFile, tourName string) {
	file, err := os.Open(sourceFile)
	if err != nil {
		fmt.Println("Error opening source file:", err)
		return
	}
	defer file.Close()

	outFile, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outFile.Close()

	scanner := bufio.NewScanner(file)
	var buffer []string
	snippets := make(map[int]TourGuide)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		buffer = append(buffer, scanner.Text())
		if len(buffer) > 8 {
			buffer = buffer[1:]
		}

		matched, err := regexp.MatchString(fmt.Sprintf(`// LivingDoc:TourGuide\("%s", (\d+), "(.*)"\)`, tourName), buffer[0])
		if err != nil {
			fmt.Println("Error matching regex:", err)
			return
		}

		if matched {
			re := regexp.MustCompile(fmt.Sprintf(`// LivingDoc:TourGuide\("%s", (\d+), "(.*)"\)`, tourName))
			matches := re.FindStringSubmatch(buffer[0])
			tourStep, _ := strconv.Atoi(matches[1])
			stepTitle := matches[2]

			// Create a new buffer to hold the lines of the code snippet without the LivingDoc tag
			var snippetBuffer []string
			for _, line := range buffer[1:] {
				if !strings.Contains(line, "LivingDoc") {
					snippetBuffer = append(snippetBuffer, line)
				}
			}

			codeSnippet := strings.Join(snippetBuffer, "\n")

			snippets[tourStep] = TourGuide{TourName: tourName, StepTitle: stepTitle, CodeSnippet: codeSnippet, LineNumber: lineNumber}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning source file:", err)
	}

	var keys []int
	for k := range snippets {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	_, err = outFile.WriteString(fmt.Sprintf("# %s\n\n", tourName))
	if err != nil {
		fmt.Println("Error writing to output file:", err)
		return
	}

	for _, k := range keys {
		tourGuide := snippets[k]
		_, err := outFile.WriteString(fmt.Sprintf("## %s\n\n```go\n%s\n```\n[Link to code](https://github.com/murex/TCR/blob/main/src/engine/./%s#L%d-L%d)\n\n", tourGuide.StepTitle, tourGuide.CodeSnippet, sourceFile, tourGuide.LineNumber, tourGuide.LineNumber+7))
		if err != nil {
			fmt.Println("Error writing to output file:", err)
			return
		}
	}
}

func main() {
	extractTourGuides("src/engine/tcr.go", "tour_guide.md", "Driver Round")
}

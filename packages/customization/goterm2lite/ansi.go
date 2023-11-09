package goterm2lite

import (
	"regexp"
	"strings"
)

// Split will follow the rules of strings.Split(input, ""), it bundles the ANSI code with other chars.
func Split(input string) []string {
	re := regexp.MustCompile(`(\x1b\[[0-9;]*m)([^\x1b]*)`)
	matches := re.FindAllStringSubmatch(input, -1)
	if len(matches) == 0 {
		return strings.Split(input, "")
	}

	output := make([]string, 0)
	for _, match := range matches {
		ansiCode := match[1]
		text := strings.Split(match[2], "")

		if len(text) > 0 {
			output = append(output, ansiCode+text[0])
		} else if len(output) >= 1 {
			output[len(output)-1] += ansiCode
		}

		/* length checks for excess */
		if len(text) == 0 {
			continue
		}

		output = append(output, text[1:]...)
	}

	return output
}

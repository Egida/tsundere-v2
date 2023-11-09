package gradient

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	Foreground = 38
	Background = 48
)

type Gradient struct {
	Colours []color
	Text    string
}

type color struct {
	Red   int
	Green int
	Blue  int
}

type escape struct {
	character string
	escape    bool
	esc       string
}

func hex2rgb(hex string) (color, error) {
	if strings.HasPrefix(hex, "#") {
		hex = hex[1:]
	}

	values, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return color{}, err
	}

	return color{
		int(values >> 16),
		int((values >> 8) & 0xFF),
		int(values & 0xFF),
	}, nil
}

// New will create the new Gradient object
func New(colours ...string) *Gradient {
	var colors []color

	for _, c := range colours {
		rgb, err := hex2rgb(c)
		if err != nil {
			return nil
		}

		colors = append(colors, rgb)
	}

	return &Gradient{
		Colours: colors,
	}
}

func (g *Gradient) Apply(ground int, text string) string {
	var output = make([]string, 0)
	var system = make([]escape, 0)

	text = strings.ReplaceAll(text, "\x1b", "\\x1b")
	var current = false

	for position := 0; position < utf8.RuneCountInString(text); position++ {
		if strings.Split(text, "")[position] == "\\" {
			var capturedCase = ""
			for attempts := position; attempts < utf8.RuneCountInString(text); attempts++ {
				position = attempts
				capturedCase += strings.Split(text, "")[attempts]
				if strings.Split(text, "")[attempts-1] == "m" && strings.Split(text, "")[attempts] != "\\" {
					break
				}
			}

			if strings.Contains(strings.Join(strings.Split(capturedCase, "")[:len(strings.Split(capturedCase, ""))-1], ""), "\\x1b[0m") {
				current = !current

				if !current {
					system = append(system, escape{character: capturedCase, escape: true})
					continue
				}
			}

			system = append(system, escape{character: capturedCase, escape: current})

			continue
		}

		system = append(system, escape{character: strings.Split(text, "")[position], escape: current})
	}

	r, gr, b := g.gradient(len(system))
	for position := range make([]string, len(system)) {
		current := system[position]

		if current.escape {
			output = append(output, fmt.Sprintf("%s%s", strings.ReplaceAll(current.esc, "\\x1b", "\x1b"), strings.ReplaceAll(current.character, "\\x1b", "\x1b")))
			continue
		}

		Red, Green, Blue := strconv.Itoa(int(r[position])), strconv.Itoa(int(gr[position])), strconv.Itoa(int(b[position])) // Converts colour types
		output = append(output, fmt.Sprintf("\x1b["+strconv.Itoa(ground)+";2;%s;%s;%sm%s\x1b[0m", Red, Green, Blue, current.character))
	}

	return strings.Join(output, "")
}

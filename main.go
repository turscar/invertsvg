package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"

	flag "github.com/spf13/pflag"
)

func main() {
	var input, output string
	flag.StringVarP(&input, "input", "i", "", "Input filename")
	flag.StringVarP(&output, "output", "o", "", "output filename")
	flag.Parse()

	content, err := os.ReadFile(input)
	if err != nil {
		log.Fatal(err)
	}
	colourRe := regexp.MustCompile(`["']#([0-9a-fA-F]{3,6})["']`)
	result := colourRe.ReplaceAllStringFunc(string(content), func(s string) string {
		var r, g, b int
		var err error
		var fixed string
		switch len(s) {
		case 9:
			_, err = fmt.Sscanf(s[1:], "#%02x%02x%02x", &r, &g, &b)
			fixed = fmt.Sprintf("%c#%02x%02x%02x%c", s[0], 255-r, 255-g, 255-b, s[0])
		case 6:
			_, err = fmt.Sscanf(s[1:], "#%01x%01x%01x", &r, &g, &b)
			fixed = fmt.Sprintf("%c#%01x%01x%01x%c", s[0], 15-r, 15-g, 15-b, s[0])
		default:
			err = errors.New("invalid colour")
			fixed = s
		}
		if err != nil {
			log.Printf("%v: %q", err, s)
			return s
		}
		return fixed
	})

	err = os.WriteFile(output, []byte(result), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

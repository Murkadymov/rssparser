package utils

import "github.com/microcosm-cc/bluemonday"

func TextCleaner(input string) string {
	cleaner := bluemonday.NewPolicy() //reg: regx, err := regexp.Compile(`<!\[CDATA\[.*?alt="(.*)]]>`) // + // (<a href=".+?>)|(</.>)|(<br>)

	cleanText := cleaner.Sanitize(input)

	return cleanText
}

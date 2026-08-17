package utils

import (
	"regexp"
	"strings"
)

var camelCaseRe = regexp.MustCompile(`([a-z0-9])([A-Z])|([A-Z]+)([A-Z][a-z])`)

func CamelToWords(s string) string {
	result := camelCaseRe.ReplaceAllString(s, "$1$3 $2$4")
	return strings.TrimSpace(result)
}

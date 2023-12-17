package utils

import (
	"fmt"
	"regexp"
)

func SequelizePlaceholders(query string) string {
	re := regexp.MustCompile(`\?`)
	placeholderCount := 1
	result := re.ReplaceAllStringFunc(query, func(match string) string {
		count := placeholderCount
		placeholderCount++
		return fmt.Sprintf("$%d", count)
	})
	return result
}

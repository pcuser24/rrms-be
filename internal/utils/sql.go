package utils

import (
	"fmt"
	"regexp"
)

// Convert the query string from the format of
// `?` to the format of `$1`, `$2`, `$3`, etc.
// This is because the `pgx` library uses the format of `$1`, `$2`, `$3`, etc.
// to represent the placeholders in the query string.
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

// Build nested SELECT query strings and arguments
// in the format SELECT ... FROM ... WHERE EXISTS (SELECT ... FROM ... WHERE EXISTS (SELECT ... FROM ... WHERE ...))
type SelectStackItem struct {
	Query string
	Args  []any
}

func BuildSelectStack(items []SelectStackItem) (string, []any) {
	return "", nil
}

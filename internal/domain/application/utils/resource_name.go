package utils

import "fmt"

func GetResourceName(aid int64) string {
	return fmt.Sprintf("[APPLICATION_%d]", aid)
}

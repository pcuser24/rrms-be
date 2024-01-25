package utils

import (
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b) // basepath is the directory of the file that calls this function
)

// GetBasePath returns the base path of the project
func GetBasePath() string {
	return basepath[0 : len(basepath)-len("/internal/utils")]
}

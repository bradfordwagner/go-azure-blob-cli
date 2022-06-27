package util

import (
	"os"
	"strings"
)

// CreateDirForFile - creates all dirs required for file to be instantiated
func CreateDirForFile(filePath string) {
	splits := strings.Split(filePath, "/")
	if len(splits) > 0 {
		prefix := strings.Join(splits[:len(splits)-1], "/")
		if _, err := os.Stat(prefix); prefix != "" && os.IsNotExist(err) {
			os.MkdirAll(prefix, 0700)
		}
	}
}

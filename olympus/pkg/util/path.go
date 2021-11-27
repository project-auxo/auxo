package util

import (
	"os"

	"github.com/project-auxo/auxo/olympus/logging"
)

var log = logging.Base()

func Validate(path string) {
	s, err := os.Stat(path)
	if err != nil {
		log.Fatalf("failed to get path: %v", err)
	}
	if s.IsDir() {
		log.Fatalf("'%s' is a directory, not a file", path)
	}
}

package main

import (
	"errors"
	"os"
)

const (
	VIDEO_DIR       = "videos"
	MAX_UPLOAD_SIZE = 1024 * 1024 * 250 // 250 MB
)

func init() {
	if _, err := os.Stat(VIDEO_DIR); errors.Is(err, os.ErrNotExist) {
		os.Mkdir(VIDEO_DIR, 0755)
	}
}

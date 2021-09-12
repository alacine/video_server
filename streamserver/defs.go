package main

import (
	"errors"
	"log"
	"os"
)

const (
	// VideoDir 视频存放目录
	VideoDir = "videos"
	// MaxUploadSize 最大上传大小
	MaxUploadSize = 1024 * 1024 * 250 // 250 MB
)

func init() {
	if _, err := os.Stat(VideoDir); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(VideoDir, 0755); err != nil {
			log.Fatalf("create storage directory failed: %s", err)
		}
	}
}

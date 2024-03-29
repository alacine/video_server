package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

func streamHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid := p.ByName("vid")
	vl := filepath.Join(VideoDir, vid)
	video, err := os.Open(vl)
	if os.IsNotExist(err) {
		sendErrorResponse(
			w,
			http.StatusNotFound,
			http.StatusText(http.StatusNotFound),
		) // 404
	} else if err != nil {
		log.Printf("Error when try to open file: %v", err)
		sendErrorResponse(
			w,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
		) // 500
	} else {
		defer func() {
			if err := video.Close(); err != nil {
				log.Printf("streamHandler close video failed %s", err)
			}
		}()
		// 告诉浏览器使用二进制流解析为 video/mp4 格式
		w.Header().Set("Content-Type", "video/mp4")
		// 二进制流传输到 Client 端
		http.ServeContent(w, r, "", time.Now(), video)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)
	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
		log.Printf("%s", err)
		sendErrorResponse(
			w,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest)+" File is too big(must <= 250MB)",
		) // 400
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		sendErrorResponse(
			w,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
		) // 500
		return
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		sendErrorResponse(
			w,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
		) // 500
	}
	vid := p.ByName("vid")
	if _, err := strconv.Atoi(vid); err != nil {
		sendErrorResponse(
			w,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest)+" Video Id must be integer",
		) // 400
		return
	}
	err = ioutil.WriteFile(filepath.Join(VideoDir, vid), data, 0666)
	if err != nil {
		log.Printf("Write file error: %v", err)
		sendErrorResponse(
			w,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
		) /// 500
		return
	}
	w.WriteHeader(http.StatusCreated) // 201
	_, err = io.WriteString(w, "Upload Successfully")
	if err != nil {
		log.Printf("sendErrorResponse error: %s", err)
	}
}

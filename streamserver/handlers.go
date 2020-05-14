package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

func streamHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid := p.ByName("vid")
	vl := VIDEO_DIR + vid
	video, err := os.Open(vl)
	defer video.Close()
	if os.IsNotExist(err) {
		sendErrorResponse(
			w,
			http.StatusNotFound,
			http.StatusText(http.StatusNotFound),
		)
	} else if err != nil {
		log.Printf("Error when try to open file: %v", err)
		sendErrorResponse(
			w,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
		)
	} else {
		// 告诉浏览器使用二进制流解析为 video/mp4 格式
		w.Header().Set("Content-Type", "video/mp4")
		// 二进制流传输到 Client 端
		http.ServeContent(w, r, "", time.Now(), video)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		log.Printf("%s", err)
		sendErrorResponse(
			w,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest)+" File is too big(must <= 250MB)",
		)
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		sendErrorResponse(
			w,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
		)
		return
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		sendErrorResponse(
			w,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
		)
	}
	vid := p.ByName("vid")
	if _, err := strconv.Atoi(vid); err != nil {
		sendErrorResponse(
			w,
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest)+" Video Id must be integer",
		)
		return
	}
	err = ioutil.WriteFile(VIDEO_DIR+vid, data, 0666)
	if err != nil {
		log.Printf("Write file error: %v", err)
		sendErrorResponse(
			w,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
		)
		return
	}
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Upload Successfully")
}

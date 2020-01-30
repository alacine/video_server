package main

import (
	"net/http"

	"github.com/alacine/video_server/scheduler/taskrunner"
	"github.com/julienschmidt/httprouter"
)

func DeleteVideoHandler() *httprouter.Router {
	router := httprouter.New()
	router.GET("/delete-video/:vid", DeleteVideo)
	return router
}

func main() {
	go taskrunner.Start()
	d := DeleteVideoHandler()
	http.ListenAndServe(":9001", d)
}

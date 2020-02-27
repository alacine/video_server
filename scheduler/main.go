package main

import (
	"net/http"

	"github.com/alacine/video_server/scheduler/taskrunner"
	"github.com/julienschmidt/httprouter"
)

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()
	router.DELETE("/scheduler/video/:vid", DeleteVideo)
	return router
}

func main() {
	go taskrunner.Start()
	r := RegisterHandlers()
	http.ListenAndServe(":9001", r)
}

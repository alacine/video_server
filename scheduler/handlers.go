package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/alacine/video_server/scheduler/dbops"
	"github.com/julienschmidt/httprouter"
)

func DeleteVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid, err := strconv.Atoi(p.ByName("vid"))
	if err != nil || vid < 1 {
		sendErrorResponse(w, http.StatusBadRequest, "video id should be int(>0)")
		return
	}
	err = dbops.AddVideoDeletionRecord(vid)
	if err != nil {
		return
	}
	log.Printf("Delete video %v", vid)
}

package main

import (
	"net/http"

	"github.com/alacine/video_server/scheduler/dbops"
	"github.com/julienschmidt/httprouter"
)

func DeleteVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid := p.ByName("vid")
	if len(vid) == 0 {
		sendErrorResponse(w, http.StatusBadRequest, "video id should not be empty")
	}
	err := dbops.AddVideoDeletionRecord(vid)
	if err != nil {
		return
	}
}

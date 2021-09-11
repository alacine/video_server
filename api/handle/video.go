package handle

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/alacine/video_server/api/dbops"
	"github.com/alacine/video_server/api/defs"
	"github.com/alacine/video_server/api/utils"
	"github.com/julienschmidt/httprouter"
)

func ListVideos(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	videos, err := dbops.ListVideos()
	if err != nil {
		log.Printf("(Error) ListVideos: %s", err)
		return
	}
	if resp, err := json.Marshal(videos); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults)
		return
	} else {
		sendNormalResponse(w, string(resp), http.StatusOK)
	}
}

func ListUserVideos(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	uid, err := strconv.Atoi(p.ByName("uid"))
	if err != nil {
		log.Printf("(Error) ListUserVideos: %s", err)
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed) // 400
		return
	}
	videos, err := dbops.ListUserVideos(uid, 0, utils.GetCurrentTimestampSec())
	if err != nil {
		log.Printf("(Error) ListUserVideos: %s", err)
		sendErrorResponse(w, defs.ErrorDBError) // 500
		return
	}
	if resp, err := json.Marshal(videos); err != nil {
		log.Printf("(Error) ListUserVideos: %s", err)
		sendErrorResponse(w, defs.ErrorInternalFaults) // 500
		return
	} else {
		sendNormalResponse(w, string(resp), http.StatusOK) // 200
	}
}

func GetVideoInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid, err := strconv.Atoi(p.ByName("vid"))
	if err != nil {
		log.Printf("(Error) GetVideoInfo: %s", err)
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed) // 400
		return
	}
	video, err := dbops.GetVideoInfo(vid)
	if err != nil {
		log.Printf("(Error) GetVideoInfo: %s", err)
		return
	}
	if resp, err := json.Marshal(video); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults) // 500
	} else {
		sendNormalResponse(w, string(resp), http.StatusOK) // 200
	}
}

func AddNewVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// https://stackoverflow.com/questions/51460418/http-request-r-formvalue-returns-nothing-map
	nvbody := &defs.NewVideo{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(nvbody); err != nil {
		log.Printf("(Error) AddNewVideo: %s", err)
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed) // 400
		return
	} else if len(nvbody.Title) == 0 {
		log.Printf("(Error) AddNewVideo: %s", err)
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed) // 400
		return
	}

	vi, err := dbops.AddNewVideo(nvbody.AuthorId, nvbody.Title, nvbody.Description)
	if err != nil {
		log.Printf("(Error) AddNewVideo: %s", err)
		sendErrorResponse(w, defs.ErrorDBError) // 500
		return
	}
	resp, err := json.Marshal(vi)
	if err != nil {
		log.Printf("(Error) AddNewVideo: %s", err)
		sendErrorResponse(w, defs.ErrorInternalFaults)
		return
	}
	sendNormalResponse(w, string(resp), http.StatusCreated) // 201
	log.Printf("AuthorId: %d, NewVideo Title: %s", nvbody.AuthorId, nvbody.Title)
}

func DeleteVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid, err := strconv.Atoi(p.ByName("vid"))
	if err != nil {
		log.Printf("(Error) DeleteVideo: %s", err)
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed) // 400
		return
	}
	err = dbops.DeleteVideoInfo(vid)
	if err != nil {
		log.Printf("(Error) DeleteVideo: %s", err)
		sendErrorResponse(w, defs.ErrorDBError) // 500
		return
	}
	go utils.SendDeleteVideoRequest(vid)
	sendNormalResponse(w, "", http.StatusNoContent) // 204
}

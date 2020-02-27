package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/alacine/video_server/api/dbops"
	"github.com/alacine/video_server/api/defs"
	"github.com/alacine/video_server/api/session"
	"github.com/alacine/video_server/api/utils"
	"github.com/julienschmidt/httprouter"
)

func CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//io.WriteString(w, "Create User Handler\n")
	res, _ := ioutil.ReadAll(r.Body)
	ubody := &defs.UserCredential{}
	if err := json.Unmarshal(res, ubody); err != nil {
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed) // 400
	}

	if err := dbops.AddUserCredential(ubody.Username, ubody.Pwd); err != nil {
		sendErrorResponse(w, defs.ErrorDBError) // 500
		return
	}

	id := session.GenerateNewSessionId(ubody.Username)
	su := &defs.SignedUp{Success: true, SessionId: id}

	if resp, err := json.Marshal(su); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults) // 500
	} else {
		sendNormalResponse(w, string(resp), http.StatusCreated) // 201
	}
}

func Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	log.Printf("%s", res)
	ubody := &defs.UserCredential{}
	if err := json.Unmarshal(res, ubody); err != nil {
		log.Printf("(Error) Login: %s", err)
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed) // 400
		return
	}
	uname := p.ByName("user_name")
	if uname != ubody.Username {
		log.Printf("(Error) Login: url's name is different with body's name")
		sendErrorResponse(w, defs.ErrorNotAuthUser) // 401
		return
	}
	pwd, err := dbops.GetUserCredential(uname)
	if err != nil || len(pwd) == 0 || pwd != ubody.Pwd {
		log.Printf("(Error) Login: user %s login failed", uname)
		sendErrorResponse(w, defs.ErrorNotAuthUser) // 401
		return
	}

	id := session.GenerateNewSessionId(uname)
	si := &defs.SignedIn{Success: true, SessionId: id}
	if resp, err := json.Marshal(si); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults) // 500
	} else {
		log.Printf("(Error) Login: user %s login succeed", uname)
		sendNormalResponse(w, string(resp), http.StatusOK)
	}
}

func GetUserInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !validateUser(r, w) {
		log.Printf("(Error) GetUserInfo: validateUser error")
		return
	}
	puname := p.ByName("user_name")
	duname, err := dbops.GetUser(puname)
	if err != nil {
		log.Printf("(Error) GetUserInfo: %s", err)
		return
	}
	ui := &defs.UserInfo{Id: duname.Id, Username: duname.Name}
	if resp, err := json.Marshal(ui); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults) // 500
	} else {
		sendNormalResponse(w, string(resp), http.StatusOK) // 200
	}
}

func ListUserVideos(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	uname := p.ByName("user_name")
	videos, err := dbops.ListUserVideos(uname, 0, utils.GetCurrentTimestampSec())
	if err != nil {
		log.Printf("(Error) ListUserVideos: %s", err)
		sendErrorResponse(w, defs.ErrorDBError) // 500
		return
	}
	vsi := &defs.VideosInfo{Videos: videos}
	if resp, err := json.Marshal(vsi); err != nil {
		log.Printf("(Error) ListUserVideos: %s", err)
		sendErrorResponse(w, defs.ErrorInternalFaults) // 500
		return
	} else {
		sendNormalResponse(w, string(resp), http.StatusOK) // 200
	}
}

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

func GetVideoInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid, err := strconv.Atoi(p.ByName("vid"))
	if err != nil {
		log.Printf("(Error) GetVideoInfo: %s", err)
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
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
	if !validateUser(r, w) {
		log.Printf("(Error) AddNewVideo: validateUser error")
		return
	}
	// https://stackoverflow.com/questions/51460418/http-request-r-formvalue-returns-nothing-map
	res, _ := ioutil.ReadAll(r.Body)
	nvbody := &defs.NewVideo{}
	if err := json.Unmarshal([]byte(res), nvbody); err != nil {
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
	if !validateUser(r, w) {
		log.Printf("(Error) DeleteVideo: validateUser error")
		return
	}
	vid, err := strconv.Atoi(p.ByName("vid"))
	if err != nil {
		log.Printf("(Error) GetVideoInfo: %s", err)
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}
	err = dbops.DeleteVideoInfo(vid)
	if err != nil {
		log.Printf("(Error) DeleteVideoInfo: %s", err)
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}
	go utils.SendDeleteVideoRequest(vid)
	sendNormalResponse(w, "", http.StatusNoContent) // 204
}

func PostComment(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !validateUser(r, w) {
		log.Printf("(Error) PostComment: validateUser error")
		return
	}
	res, _ := ioutil.ReadAll(r.Body)
	cbody := &defs.NewComment{}
	if err := json.Unmarshal(res, cbody); err != nil {
		log.Printf("(Error) PostComment: %s", err)
		sendErrorResponse(w, defs.ErrorInternalFaults) // 500
		return
	}
	vid := p.ByName("vid")
	if err := dbops.AddNewComment(vid, cbody.AuthorId, cbody.Content); err != nil {
		log.Printf("(Error) PostComment: %s", err)
		sendErrorResponse(w, defs.ErrorDBError) // 500
	} else {
		sendNormalResponse(w, "ok", http.StatusCreated) // 201
	}
}

func ListComments(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid, err := strconv.Atoi(p.ByName("vid"))
	if err != nil {
		log.Printf("(Error) GetVideoInfo: %s", err)
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
		return
	}
	comments, err := dbops.ListComments(vid, 0, utils.GetCurrentTimestampSec())
	if err != nil {
		log.Printf("(Error) ListComments: %s", err)
		sendErrorResponse(w, defs.ErrorDBError) // 500
		return
	}
	cs := &defs.Comments{Comment: comments}
	if resp, err := json.Marshal(cs); err != nil {
		log.Printf("(Error) ListComments: %s", err)
		sendErrorResponse(w, defs.ErrorInternalFaults) // 500
		return
	} else {
		sendNormalResponse(w, string(resp), http.StatusOK)
	}
}
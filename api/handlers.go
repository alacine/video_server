package main

import (
	"encoding/json"
	"io"
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
	ubody := &defs.UserCredential{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(ubody); err != nil && err != io.EOF {
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed) // 400
		log.Printf("(Error) CreateUser: %s", err)
		return
	}
	uid, err := dbops.AddUserCredential(ubody.Username, ubody.Pwd)
	if err != nil {
		sendErrorResponse(w, defs.ErrorDBError) // 500
		log.Printf("(Error) CreateUser: %s", err)
		return
	}
	new_user := &defs.UserInfo{
		Username: ubody.Username,
		Id:       uid,
	}
	resp, err := json.Marshal(new_user)
	if err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults) // 500
		log.Printf("(Error) CreateUser: %s", err)
		return
	}
	sendNormalResponse(w, string(resp), http.StatusCreated) // 201
}

func GetUserInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//if !validateUser(r, w) {
	//log.Printf("(Error) GetUserInfo: validateUser error")
	//return
	//}
	uid, err := strconv.Atoi(p.ByName("uid"))
	if err != nil {
		log.Printf("(Error) GetUserInfo: %s", err)
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed) // 400
		return
	}
	user, err := dbops.GetUser(uid)
	if err != nil {
		log.Printf("(Error) GetUserInfo: %s", err)
		return
	}
	ui := &defs.UserInfo{Id: user.Id, Username: user.Name}
	if resp, err := json.Marshal(ui); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults) // 500
	} else {
		sendNormalResponse(w, string(resp), http.StatusOK) // 200
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

func Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ubody := &defs.UserCredential{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(ubody); err != nil {
		log.Printf("(Error) Login: %s", err)
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed) // 400
		return
	}
	log.Printf("%#v", ubody)
	uid, pwd, err := dbops.GetUserCredential(ubody.Username)
	if err != nil || len(pwd) == 0 || pwd != ubody.Pwd {
		log.Printf("(Error) Login: user %s login failed", ubody.Username)
		sendErrorResponse(w, defs.ErrorNotAuthUser) // 401
		return
	}

	id := session.GenerateNewSessionId(uid)
	si := &defs.SignedIn{Success: true, SessionId: id, UserId: uid}
	if resp, err := json.Marshal(si); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults) // 500
	} else {
		log.Printf("Login: user %s login succeed", ubody.Username)
		sendNormalResponse(w, string(resp), http.StatusOK) /// 200
	}
}

func Logout(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !validateUser(r, w) {
		log.Printf("(Error) Logout: validateUser error")
		return
	}
	sid, err := r.Cookie(HEADER_FIELD_SESSION)
	if err != nil {
		sendErrorResponse(w, defs.ErrorNotAuthUser) // 401
	}
	session.DeleteExpiredSession(sid.Value)
	sendNormalResponse(w, "Logout", http.StatusResetContent) // 205
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
	if !validateUser(r, w) {
		log.Printf("(Error) AddNewVideo: validateUser error")
		return
	}
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
	if !validateUser(r, w) {
		log.Printf("(Error) DeleteVideo: validateUser error")
		return
	}
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

func PostComment(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !validateUser(r, w) {
		log.Printf("(Error) PostComment: validateUser error")
		return
	}
	cbody := &defs.NewComment{}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(cbody); err != nil {
		log.Printf("(Error) PostComment: %s", err)
		sendErrorResponse(w, defs.ErrorInternalFaults) // 500
		return
	}
	vid, err := strconv.Atoi(p.ByName("vid"))
	if err != nil {
		log.Printf("(Error) PostComment: %s", err)
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed) // 400
	}
	log.Printf("postcomment %v, %v", vid, cbody)
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
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed) /// 400
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
		sendNormalResponse(w, string(resp), http.StatusOK) /// 200
	}
}

package main

import (
	"encoding/json"
	"fmt"
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
	ui := &defs.UserInfo{Id: duname.Id, Username: duname.LoginName}
	if resp, err := json.Marshal(ui); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults) // 500
	} else {
		sendNormalResponse(w, string(resp), http.StatusOK) // 200
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
	res := fmt.Sprintf("{\"author_id\":%s,\"name\":\"%s\"}", r.FormValue("author_id"), r.FormValue("name"))
	//res, _ := ioutil.ReadAll(r.Body)
	nvbody := &defs.NewVideo{}
	if err := json.Unmarshal([]byte(res), nvbody); err != nil {
		log.Printf("(Error) AddNewVideo: %s", err)
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed) // 400
		return
	}

	vi, err := dbops.AddNewVideo(nvbody.AuthorId, nvbody.Name)
	log.Printf("AuthorId: %d, NewVideo name: %s", nvbody.AuthorId, nvbody.Name)
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

	// save video
	r.Body = http.MaxBytesReader(w, r.Body, defs.MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(defs.MAX_UPLOAD_SIZE); err != nil {
		dbops.DeleteVideoInfo(vi.Id)
		sendErrorResponse(w, defs.ErrorFileSize)
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		log.Printf("(Error) UploadVideo: %s", err)
		dbops.DeleteVideoInfo(vi.Id)
		sendErrorResponse(w, defs.ErrorInternalFaults)
		return
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("(Error) UploadVideo: %s", err)
		dbops.DeleteVideoInfo(vi.Id)
		sendErrorResponse(w, defs.ErrorInternalFaults)
	}
	log.Printf(defs.VIDEO_DIR + strconv.Itoa(vi.Id))
	err = ioutil.WriteFile(defs.VIDEO_DIR+strconv.Itoa(vi.Id), data, 0666)
	if err != nil {
		log.Printf("(Error) UploadVideo: %s", err)
		dbops.DeleteVideoInfo(vi.Id)
		sendErrorResponse(w, defs.ErrorInternalFaults)
		return
	}
	sendNormalResponse(w, string(resp), http.StatusCreated) // 201
}

func ListUserVideos(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	uname := p.ByName("user_name")
	videos, err := dbops.ListVideoInfo(uname, 0, utils.GetCurrentTimestampSec())
	if err != nil {
		log.Printf("(Error) ListAllVideos: %s", err)
		sendErrorResponse(w, defs.ErrorDBError) // 500
		return
	}
	vsi := &defs.VideosInfo{Videos: videos}
	if resp, err := json.Marshal(vsi); err != nil {
		log.Printf("(Error) ListAllVideos: %s", err)
		sendErrorResponse(w, defs.ErrorInternalFaults) // 500
		return
	} else {
		sendNormalResponse(w, string(resp), http.StatusOK) // 200
	}
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
	// TODO
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

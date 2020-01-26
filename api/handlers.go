package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/alacine/video_server/api/dbops"
	"github.com/alacine/video_server/api/defs"
	"github.com/alacine/video_server/api/session"
	"github.com/julienschmidt/httprouter"
)

func CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//io.WriteString(w, "Create User Handler\n")
	res, _ := ioutil.ReadAll(r.Body)
	ubody := &defs.UserCredential{}
	if err := json.Unmarshal(res, ubody); err != nil {
		sendErrorResponse(w, defs.ErrorRequestBodyParseFailed)
	}

	if err := dbops.AddUserCredential(ubody.Username, ubody.Pwd); err != nil {
		sendErrorResponse(w, defs.ErrorDBError)
		return
	}

	id := session.GenerateNewSessionId(ubody.Username)
	su := &defs.SignedUp{Success: true, SessionId: id}

	if resp, err := json.Marshal(su); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults)
	} else {
		sendNormalResponse(w, string(resp), 201)
	}
}

func Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	uname := p.ByName("user_name")
	io.WriteString(w, uname)
}

package handle

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/alacine/video_server/api/dbops"
	"github.com/alacine/video_server/api/defs"
	"github.com/julienschmidt/httprouter"
)

// CreateUser ...
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
	user := &defs.UserInfo{
		Username: ubody.Username,
		ID:       uid,
	}
	resp, err := json.Marshal(user)
	if err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults) // 500
		log.Printf("(Error) CreateUser: %s", err)
		return
	}
	sendNormalResponse(w, string(resp), http.StatusCreated) // 201
}

// GetUserInfo ...
func GetUserInfo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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
	ui := &defs.UserInfo{ID: user.ID, Username: user.Name}
	if resp, err := json.Marshal(ui); err != nil {
		sendErrorResponse(w, defs.ErrorInternalFaults) // 500
	} else {
		sendNormalResponse(w, string(resp), http.StatusOK) // 200
	}
}

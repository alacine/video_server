package handle

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/alacine/video_server/api/dbops"
	"github.com/alacine/video_server/api/defs"
	"github.com/alacine/video_server/api/middleware"
	"github.com/alacine/video_server/api/session"
	"github.com/julienschmidt/httprouter"
)

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
	sid, err := r.Cookie(middleware.HEADER_FIELD_SESSION)
	if err != nil {
		sendErrorResponse(w, defs.ErrorNotAuthUser) // 401
	}
	session.DeleteExpiredSession(sid.Value)
	sendNormalResponse(w, "Logout", http.StatusResetContent) // 205
}

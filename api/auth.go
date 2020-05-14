package main

import (
	"net/http"
	"strconv"

	"github.com/alacine/video_server/api/defs"
	"github.com/alacine/video_server/api/session"
)

var HEADER_FIELD_SESSION = "X-Session-Id"
var HEADER_FIELD_UID = "X-User-Id"

// user 校验
func validateUser(r *http.Request, w http.ResponseWriter) bool {
	session_id, err := r.Cookie(HEADER_FIELD_SESSION)
	if err != nil || len(session_id.Value) == 0 {
		sendErrorResponse(w, defs.ErrorNotAuthUser) // 401
		return false
	}
	uidstr, err1 := r.Cookie(HEADER_FIELD_UID)
	uid, err2 := strconv.Atoi(uidstr.Value)
	if err1 != nil || err2 != nil {
		sendErrorResponse(w, defs.ErrorNotAuthUser) // 401
		return false
	}
	if suid, ok := session.IsSessionExpired(session_id.Value); ok == true || suid != uid {
		sendErrorResponse(w, defs.ErrorNotAuthUser) // 401
		return false
	}
	return true
}

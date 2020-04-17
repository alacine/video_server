package main

import (
	"net/http"

	"github.com/alacine/video_server/api/defs"
	"github.com/alacine/video_server/api/session"
)

var HEADER_FIELD_SESSION = "X-Session-Id"
var HEADER_FIELD_UNAME = "X-User-Name"

// session 校验
func validateUserSession(r *http.Request) bool {
	sid := r.Header.Get(HEADER_FIELD_SESSION)
	if len(sid) == 0 {
		return false
	}
	uname, ok := session.IsSessionExpired(sid)
	if ok {
		return false
	}
	r.Header.Add(HEADER_FIELD_UNAME, uname)
	return true
}

// user 校验
func validateUser(r *http.Request, w http.ResponseWriter) bool {
	session_id, err := r.Cookie(HEADER_FIELD_SESSION)
	if err != nil {
		return false
	}
	uname, err := r.Cookie(HEADER_FIELD_UNAME)
	if err != nil {
		return false
	}
	if len(uname.Value) == 0 || len(session_id.Value) == 0 {
		sendErrorResponse(w, defs.ErrorNotAuthUser) // 401
		return false
	}
	return true
}

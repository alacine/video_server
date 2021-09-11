package middleware

import (
	"net/http"
	"strconv"

	"github.com/alacine/video_server/api/session"
	"github.com/julienschmidt/httprouter"
)

var HEADER_FIELD_SESSION = "X-Session-Id"
var HEADER_FIELD_UID = "X-User-Id"

// 检查登录状态的中间件
func CheckLogin(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if !validateUser(r) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next(w, r, ps)
	}
}

// user 校验
func validateUser(r *http.Request) bool {
	session_id, err := r.Cookie(HEADER_FIELD_SESSION)
	if err != nil || len(session_id.Value) == 0 {
		//sendErrorResponse(w, defs.ErrorNotAuthUser) // 401
		return false
	}
	uidstr, err1 := r.Cookie(HEADER_FIELD_UID)
	uid, err2 := strconv.Atoi(uidstr.Value)
	if err1 != nil || err2 != nil {
		//sendErrorResponse(w, defs.ErrorNotAuthUser) // 401
		return false
	}
	if suid, ok := session.IsSessionExpired(session_id.Value); ok == true || suid != uid {
		//sendErrorResponse(w, defs.ErrorNotAuthUser) // 401
		return false
	}
	return true
}

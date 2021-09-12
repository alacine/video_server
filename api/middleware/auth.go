package middleware

import (
	"net/http"
	"strconv"

	"github.com/alacine/video_server/api/session"
	"github.com/julienschmidt/httprouter"
)

const (
	// HeaderFieldSession ...
	HeaderFieldSession = "X-Session-Id"
	// HeaderFieldUID ...
	HeaderFieldUID = "X-User-Id"
)

// CheckLogin 检查登录状态的中间件
func CheckLogin(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if !validateUser(r) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next(w, r, ps)
	}
}

// validateUser 用户校验
func validateUser(r *http.Request) bool {
	sid, err := r.Cookie(HeaderFieldSession)
	if err != nil || len(sid.Value) == 0 {
		//sendErrorResponse(w, defs.ErrorNotAuthUser) // 401
		return false
	}
	uidstr, err1 := r.Cookie(HeaderFieldUID)
	uid, err2 := strconv.Atoi(uidstr.Value)
	if err1 != nil || err2 != nil {
		//sendErrorResponse(w, defs.ErrorNotAuthUser) // 401
		return false
	}
	if suid, ok := session.IsSessionExpired(sid.Value); ok || suid != uid {
		//sendErrorResponse(w, defs.ErrorNotAuthUser) // 401
		return false
	}
	return true
}

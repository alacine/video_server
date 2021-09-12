package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type middleWareHandler struct {
	r *httprouter.Router
	l *ConnLimiter
}

func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !m.l.GetConnLimiter() {
		sendErrorResponse(w, http.StatusTooManyRequests, "Too many requests")
		return
	}
	defer m.l.ReleaseConn()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, DELETE, PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Session-Id, X-User-Id, Cookie")
	w.Header().Add("Set-Cookie", "HttpOnly;Secure;SameSite=Strict")
	m.r.ServeHTTP(w, r)
}

// NewMiddleWareHandler ...
func NewMiddleWareHandler(r *httprouter.Router, cc int) http.Handler {
	m := middleWareHandler{}
	m.r = r
	m.l = NewConnLimiter(cc)
	return m
}

// RegisterHandlers ...
func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()
	router.GET("/stream/videos/:vid", streamHandler)
	router.POST("/stream/videos/:vid", uploadHandler)
	return router
}

func main() {
	r := RegisterHandlers()
	mh := NewMiddleWareHandler(r, 5)
	log.Fatalln(http.ListenAndServe(":9000", mh))
}

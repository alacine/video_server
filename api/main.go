package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type middleWareHandler struct {
	r *httprouter.Router
}

func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// check session
	//validateUserSession(r)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, DELETE")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, X-Session-Id, X-User-Id, Cookie")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Set-Cookie", "HttpOnly;Secure;SameSite=Strict")
	m.r.ServeHTTP(w, r)
}

func NewMiddleWareHandler(r *httprouter.Router, c int) http.Handler {
	m := middleWareHandler{}
	m.r = r
	return m
}

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.POST("/api/users", CreateUser)
	router.GET("/api/users/:uid", GetUserInfo)
	router.GET("/api/users/:uid/videos", ListUserVideos)

	router.POST("/api/sessions", Login)
	router.DELETE("/api/sessions", Logout)

	router.GET("/api/videos", ListVideos)
	router.GET("/api/videos/:vid", GetVideoInfo)
	router.POST("/api/videos", AddNewVideo)
	router.DELETE("/api/videos/:vid", DeleteVideo)

	router.POST("/api/videos/:vid/comments", PostComment)
	router.GET("/api/videos/:vid/comments", ListComments)

	return router
}

func main() {
	r := RegisterHandlers()
	mh := NewMiddleWareHandler(r, 3)
	http.ListenAndServe(":8000", mh)
}

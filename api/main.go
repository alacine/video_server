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
	validateUserSession(r)
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, DELETE")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, X-Session-Id, X-X-User-Name, Cookie")
	w.Header().Set("Content-Type", "application/json")
	r.Header.Set("Set-Cookie", "HttpOnly;Secure;SameSite=Strict")
	m.r.ServeHTTP(w, r)
}

func NewMiddleWareHandler(r *httprouter.Router, c int) http.Handler {
	m := middleWareHandler{}
	m.r = r
	return m
}

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.POST("/api/user", CreateUser)
	router.POST("/api/user/:user_name", Login)
	router.GET("/api/user/:user_name", GetUserInfo)
	router.GET("/api/user/:user_name/videos", ListUserVideos)

	//router.POST("/user/:user_name/video", AddNewVideo)
	//router.GET("/video/:vid", StreamVideo)
	router.GET("/api/videos", ListVideos)
	router.GET("/api/video/:vid", GetVideoInfo)
	router.POST("/api/video", AddNewVideo)
	router.DELETE("/api/video/:vid", DeleteVideo)

	router.POST("/api/video/:vid/comments", PostComment)
	router.GET("/api/video/:vid/comments", ListComments)

	return router
}

func main() {
	r := RegisterHandlers()
	mh := NewMiddleWareHandler(r, 3)
	http.ListenAndServe(":8000", mh)
}

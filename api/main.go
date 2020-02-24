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
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
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

	router.POST("/user", CreateUser)
	router.POST("/user/:user_name", Login)
	router.GET("/user/:user_name", GetUserInfo)
	router.GET("/user/:user_name/videos", ListUserVideos)

	//router.POST("/user/:user_name/video", AddNewVideo)
	//router.GET("/video/:vid", StreamVideo)
	router.GET("/video/:vid/info", GetVideoInfo)
	router.POST("/video", AddNewVideo)
	router.DELETE("/videos/:vid/delete", DeleteVideo)

	router.POST("/video/:vid/comments", PostComment)
	router.GET("/video/:vid/comments", ListComments)

	return router
}

func main() {
	r := RegisterHandlers()
	mh := NewMiddleWareHandler(r, 3)
	http.ListenAndServe(":8000", mh)
}

package main

import (
	"log"
	"net/http"

	"github.com/alacine/video_server/api/handle"
	"github.com/alacine/video_server/api/middleware"
	"github.com/alacine/video_server/api/session"
	"github.com/julienschmidt/httprouter"
)

type middleWareHandler struct {
	r *httprouter.Router
}

func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// cors
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, DELETE, PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Session-ID, X-User-ID, Cookie")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Session-ID, X-User-ID, Cookie")
	w.Header().Add("Set-Cookie", "HttpOnly;Secure;SameSite=Strict")
	// json header
	w.Header().Set("Content-Type", "application/json")
	m.r.ServeHTTP(w, r)
}

// NewMiddleWareHandler 用来放全局中间件
func NewMiddleWareHandler(r *httprouter.Router) http.Handler {
	m := middleWareHandler{}
	m.r = r
	return m
}

// Middleware 中间件
type Middleware func(httprouter.Handle) httprouter.Handle

// Inline 包装单独的中间件
func Inline(origin httprouter.Handle, mws ...Middleware) httprouter.Handle {
	for i := len(mws) - 1; i >= 0; i-- {
		origin = mws[i](origin)
	}
	return origin
}

// RegisterHandlers 注册路由
func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()

	router.POST("/api/users", handle.CreateUser)
	router.GET("/api/users/:uid", handle.GetUserInfo)
	router.GET("/api/users/:uid/videos", handle.ListUserVideos)

	router.POST("/api/sessions", handle.Login)
	router.DELETE("/api/sessions", Inline(handle.Logout, middleware.CheckLogin))

	router.GET("/api/videos", handle.ListVideos)
	router.GET("/api/videos/:vid", handle.GetVideoInfo)
	router.POST("/api/videos", Inline(handle.AddNewVideo, middleware.CheckLogin))
	router.DELETE("/api/videos/:vid", handle.DeleteVideo)

	router.POST("/api/videos/:vid/comments", Inline(handle.PostComment, middleware.CheckLogin))
	router.GET("/api/videos/:vid/comments", handle.ListComments)

	return router
}

// Prepare 载入数据库中原有 session
func Prepare() {
	session.LoadSessionsFromDB()
}

func main() {
	Prepare()
	r := RegisterHandlers()
	mh := NewMiddleWareHandler(r)
	log.Fatalln(http.ListenAndServe(":8000", mh))
}

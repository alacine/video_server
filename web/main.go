package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func RegisterHandlers() *httprouter.Router {
	router := httprouter.New()
	router.GET("/", homeHandler)
	router.POST("/", homeHandler)
	router.GET("/userhome", userHomeHandler)
	router.POST("/userhome", userHomeHandler)
	router.POST("/api", apiHandler)
	router.ServeFiles("/statics/*filepath", http.Dir("./template"))
	return router
}

func main() {
	r := RegisterHandlers()
	http.ListenAndServe(":8080", r)
}

package main

import (
	"io"
	"log"
	"net/http"
)

func sendErrorResponse(w http.ResponseWriter, sc int, errMsg string) {
	w.WriteHeader(sc)
	_, err := io.WriteString(w, errMsg)
	if err != nil {
		log.Printf("sendErrorResponse error: %s", err)
	}
}

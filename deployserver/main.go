package main

import (
	"io"
	"log"
	"net/http"
	"os/exec"
)

func reLaunch() {
	cmd := exec.Command("bash", "./deploy.sh")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Wait()
}

func showPage(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<h1> Hello, this is video_server's deployee page </h1>")
	reLaunch()
}

func main() {
	http.HandleFunc("/", showPage)
	http.ListenAndServe(":8001", nil)
}

package handle

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/alacine/video_server/api/defs"
)

func sendErrorResponse(w http.ResponseWriter, errResp defs.ErrResponse) {
	w.WriteHeader(errResp.HTTPSC)
	resStr, _ := json.Marshal(&errResp.Error)
	_, err := io.WriteString(w, string(resStr))
	if err != nil {
		log.Printf("sendErrorResponse error: %s", err)
	}
}

func sendNormalResponse(w http.ResponseWriter, resp string, sc int) {
	w.WriteHeader(sc)
	_, err := io.WriteString(w, resp)
	if err != nil {
		log.Printf("sendNormalResponse error: %s", err)
	}
}

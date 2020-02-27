package utils

import (
	"alacine/config"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func NewUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func GetCurrentTimestampSec() int {
	timestamp, _ := strconv.Atoi(strconv.FormatInt(time.Now().UnixNano()/1e9, 10))
	return timestamp
}

func SendDeleteVideoRequest(vid int) {
	addr := config.GetScheduler()
	url := "http://" + addr + "/video/" + strconv.Itoa(vid)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Printf("(Error) SendDeleteVideoRequest error")
	}
	// https://stackoverflow.com/questions/46310113/consume-a-delete-endpoint-from-golang
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Printf("(Error) SendDeleteVideoRequest error")
	}
}

package taskrunner

import (
	"errors"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/alacine/video_server/scheduler/dbops"
)

func deleteVideo(vid int) error {
	err := os.Remove(VIDEO_DIR + strconv.Itoa(vid))
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Deleting video error")
		return err
	}
	return nil
}

func VideoClearDispatcher(dc dataChan) error {
	res, err := dbops.ReadVideoDeletionRecord(3)
	if err != nil {
		log.Printf("Video clear dispatcher error: %v", err)
		return err
	}
	if len(res) == 0 {
		return errors.New("All tasks finished")
	}
	for _, id := range res {
		dc <- id
	}
	return nil
}

func VideoClearExecutor(dc dataChan) error {
	errMap := sync.Map{}
forloop:
	for {
		select {
		case vid := <-dc:
			// 这里可能会有重复读写，但是不影响最终的结果
			go func(id interface{}) {
				if err := deleteVideo(id.(int)); err != nil {
					errMap.Store(id, err)
				}
				if err := dbops.DelVideoDeletionRecord(id.(int)); err != nil {
					return
				}
			}(vid)
		default:
			break forloop
		}
	}
	errMap.Range(func(k, v interface{}) bool {
		if v.(error) != nil {
			return false
		}
		return true
	})
	return nil
}

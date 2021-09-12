package session

import (
	"log"
	"sync"
	"time"

	"github.com/alacine/video_server/api/dbops"
	"github.com/alacine/video_server/api/defs"
	"github.com/alacine/video_server/api/utils"
)

var sessionMap *sync.Map

func init() {
	sessionMap = &sync.Map{}
}

func nowInMilli() int64 {
	return time.Now().UnixNano() / 1e6
}

// DeleteExpiredSession 删除会话信息
func DeleteExpiredSession(sid string) {
	sessionMap.Delete(sid)
	err := dbops.DeleteSession(sid)
	if err != nil {
		log.Println(err)
	}
}

// LoadSessionsFromDB 加载会话信息
func LoadSessionsFromDB() {
	r, err := dbops.RetrieveAllSessions()
	if err != nil {
		return
	}
	r.Range(func(k, v interface{}) bool {
		ss := v.(*defs.SimpleSession)
		sessionMap.Store(k, ss)
		return true
	})
}

// GenerateNewSessionID 生成会话信息
func GenerateNewSessionID(uid int) string {
	id, _ := utils.NewUUID()
	ct := nowInMilli()
	ttl := ct + 300*60*1000 // Serverside session vaild time: 300 min
	ss := &defs.SimpleSession{UserID: uid, TTL: ttl}
	sessionMap.Store(id, ss)
	err := dbops.InsertSession(id, ttl, uid)
	if err != nil {
		log.Printf("GenerateNewSessionID insert db failed: %s", err)
	}
	return id
}

// IsSessionExpired 判断会话信息是否过期
func IsSessionExpired(sid string) (int, bool) {
	ss, ok := sessionMap.Load(sid)
	if ok {
		ct := nowInMilli()
		if ss.(*defs.SimpleSession).TTL < ct {
			DeleteExpiredSession(sid)
			return 0, true
		}
		return ss.(*defs.SimpleSession).UserID, false
	}
	return 0, true
}

package session

import (
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

func DeleteExpiredSession(sid string) {
	sessionMap.Delete(sid)
	dbops.DeleteSession(sid)
}

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

func GenerateNewSessionId(uid int) string {
	id, _ := utils.NewUUID()
	ct := nowInMilli()
	ttl := ct + 300*60*1000 // Serverside session vaild time: 300 min
	ss := &defs.SimpleSession{UserId: uid, TTL: ttl}
	sessionMap.Store(id, ss)
	dbops.InsertSession(id, ttl, uid)
	return id
}

func IsSessionExpired(sid string) (int, bool) {
	ss, ok := sessionMap.Load(sid)
	if ok {
		ct := nowInMilli()
		if ss.(*defs.SimpleSession).TTL < ct {
			DeleteExpiredSession(sid)
			return 0, true
		}
		return ss.(*defs.SimpleSession).UserId, false
	}
	return 0, true
}

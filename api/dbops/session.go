package dbops

import (
	"database/sql"
	"log"
	"strconv"
	"sync"

	"github.com/alacine/video_server/api/defs"
)

// InsertSession ...
func InsertSession(sid string, ttl int64, uid int) error {
	ttlstr := strconv.FormatInt(ttl, 10)
	stmt, err := dbConn.Prepare(`INSERT INTO sessions (session_id, TTL, uid)
									VALUES (?, ?, ?)`)
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("close db connection failed: %s", err)
		}
	}()
	if err != nil {
		return err
	}
	_, err = stmt.Exec(sid, ttlstr, uid)
	if err != nil {
		return err
	}
	return nil
}

// RetrieveSession ...
func RetrieveSession(sid string) (*defs.SimpleSession, error) {
	ss := &defs.SimpleSession{}
	stmt, err := dbConn.Prepare(`SELECT TTL, name
									FROM sessions WHERE session_id = ?`)
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("close db connection failed: %s", err)
		}
	}()
	if err != nil {
		return nil, err
	}
	var ttl string
	var uid int
	err = stmt.QueryRow(sid).Scan(&ttl, &uid)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if res, err := strconv.ParseInt(ttl, 10, 64); err == nil {
		ss.TTL = res
		ss.UserID = uid
		ss.UserID = uid
	} else {
		return nil, err
	}
	return ss, nil
}

// RetrieveAllSessions 当服务重启的时候需要进行此操作, 重新获取数据库中所有的 session 到缓存中
func RetrieveAllSessions() (*sync.Map, error) {
	m := &sync.Map{}
	stmt, err := dbConn.Prepare("SELECT * FROM sessions")
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("close db connection failed: %s", err)
		}
	}()
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}
	rows, err := stmt.Query()
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}
	for rows.Next() {
		var id string
		var ttlstr string
		var uid int
		if err := rows.Scan(&id, &ttlstr, &uid); err != nil {
			log.Printf("retrieve sessions error: %s", err)
		}
		if ttl, err := strconv.ParseInt(ttlstr, 10, 64); err == nil {
			ss := &defs.SimpleSession{UserID: uid, TTL: ttl}
			m.Store(id, ss)
			log.Printf(" session id: %s, ttl: %d", id, ss.TTL)
		} else {
			return nil, err
		}
	}
	return m, nil
}

// DeleteSession ...
func DeleteSession(sid string) error {
	stmt, err := dbConn.Prepare("delete from sessions where session_id = ?")
	if err != nil {
		log.Printf("%s", err)
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("close db connection failed: %s", err)
		}
	}()
	if _, err := stmt.Exec(sid); err != nil {
		return err
	}
	return nil
}

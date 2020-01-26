package dbops

import (
	"database/sql"
	"log"
	"strconv"
	"sync"

	"github.com/alacine/video_server/api/defs"
)

func InsertSession(sid string, ttl int64, uname string) error {
	ttlstr := strconv.FormatInt(ttl, 10)
	stmtIns, err := dbConn.Prepare(`INSERT INTO sessions (session_id, TTL, login_name) 
									VALUES (?, ?, ?)`)
	defer stmtIns.Close()
	if err != nil {
		return err
	}
	_, err = stmtIns.Exec(sid, ttlstr, uname)
	if err != nil {
		return err
	}
	return nil
}

func RetrieveSession(sid string) (*defs.SimpleSession, error) {
	ss := &defs.SimpleSession{}
	stmtOut, err := dbConn.Prepare(`SELECT TTL, login_name 
									FROM sessions WHERE session_id = ?`)
	defer stmtOut.Close()
	if err != nil {
		return nil, err
	}
	var ttl string
	var uname string
	err = stmtOut.QueryRow(sid).Scan(&ttl, &uname)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if res, err := strconv.ParseInt(ttl, 10, 64); err == nil {
		ss.TTL = res
		ss.Username = uname
	} else {
		return nil, err
	}
	return ss, nil
}

// 当服务重启的时候需要进行此操作, 重新获取数据库中所有的 session 到缓存中
func RetrieveAllSessions() (*sync.Map, error) {
	m := &sync.Map{}
	stmtOut, err := dbConn.Prepare("SELECT * FROM sessions")
	defer stmtOut.Close()
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}
	rows, err := stmtOut.Query()
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}
	for rows.Next() {
		var id string
		var ttlstr string
		var login_name string
		if err := rows.Scan(&id, ttlstr, login_name); err != nil {
			log.Printf("retrieve sessions error: %s", err)
		}
		if ttl, err := strconv.ParseInt(ttlstr, 10, 64); err == nil {
			ss := &defs.SimpleSession{Username: login_name, TTL: ttl}
			m.Store(id, ss)
			log.Printf(" session id: %s, ttl: %d", id, ss.TTL)
		} else {
			return nil, err
		}
	}
	return m, nil
}

func DeleteSession(sid string) error {
	stmtDel, err := dbConn.Prepare("delete from sessions where id = ?")
	if err != nil {
		log.Printf("%s", err)
		return err
	}
	defer stmtDel.Close()
	if _, err := stmtDel.Exec(sid); err != nil {
		return err
	}
	return nil
}

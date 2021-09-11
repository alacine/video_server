package dbops

import (
	"database/sql"
	"log"

	"github.com/alacine/video_server/api/defs"
	_ "github.com/go-sql-driver/mysql"
)

func AddUserCredential(name string, pwd string) (int, error) {
	stmtIns, err := dbConn.Prepare("INSERT INTO users (name, pwd) VALUES (?, ?)")
	defer stmtIns.Close()
	if err != nil {
		log.Printf("(ERROR) AddUserCredential sql prepare error: %s", err)
		return -1, err
	}
	result, err := stmtIns.Exec(name, pwd)
	if err != nil {
		log.Printf("(ERROR) AddUserCredential sql exec error: %s", err)
		return -1, err
	}
	vid, err := result.LastInsertId()
	if err != nil {
		log.Printf("(ERROR) AddUserCredential sql exec error: %s", err)
		return -1, err
	}
	return int(vid), nil
}

func GetUserCredential(name string) (int, string, error) {
	stmtOut, err := dbConn.Prepare("SELECT id, pwd FROM users WHERE name = ?")
	defer stmtOut.Close()
	if err != nil {
		log.Printf("(ERROR) GetUserCredential sql prepare error %s", err)
		return 0, "", err
	}
	var pwd string
	var uid int
	err = stmtOut.QueryRow(name).Scan(&uid, &pwd)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("(ERROR) GetUserCredential sql query error %s", err)
		return 0, "", err
	}
	return uid, pwd, nil
}

func GetUser(uid int) (*defs.User, error) {
	stmtOut, err := dbConn.Prepare("select name from users where id = ?")
	defer stmtOut.Close()
	if err != nil {
		log.Printf("(ERROR) GetUser sql prepare error: %s", err)
		return nil, err
	}
	var uname string
	err = stmtOut.QueryRow(uid).Scan(&uname)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("(ERROR) GetUser sql query error: %s", err)
		return nil, err
	}
	res := &defs.User{Id: uid, Name: uname}
	return res, nil
}

func DeleteUser(name string, pwd string) error {
	stmtDel, err := dbConn.Prepare("DELETE FROM users WHERE name = ? AND pwd = ?")
	defer stmtDel.Close()
	if err != nil {
		log.Printf("(ERROR) DeleteUser sql prepare error: %s", err)
		return err
	}
	_, err = stmtDel.Exec(name, pwd)
	if err != nil {
		log.Printf("(ERROR) DeleteUser sql exec error: %s", err)
		return err
	}
	return nil
}

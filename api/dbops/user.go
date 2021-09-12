package dbops

import (
	"database/sql"
	"log"

	"github.com/alacine/video_server/api/defs"
	// Registe mysql
	_ "github.com/go-sql-driver/mysql"
)

// AddUserCredential ...
func AddUserCredential(name string, pwd string) (int, error) {
	stmt, err := dbConn.Prepare("INSERT INTO users (name, pwd) VALUES (?, ?)")
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("close db connection failed: %s", err)
		}
	}()
	if err != nil {
		log.Printf("(ERROR) AddUserCredential sql prepare error: %s", err)
		return -1, err
	}
	result, err := stmt.Exec(name, pwd)
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

// GetUserCredential ...
func GetUserCredential(name string) (int, string, error) {
	stmt, err := dbConn.Prepare("SELECT id, pwd FROM users WHERE name = ?")
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("close db connection failed: %s", err)
		}
	}()
	if err != nil {
		log.Printf("(ERROR) GetUserCredential sql prepare error %s", err)
		return 0, "", err
	}
	var pwd string
	var uid int
	err = stmt.QueryRow(name).Scan(&uid, &pwd)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("(ERROR) GetUserCredential sql query error %s", err)
		return 0, "", err
	}
	return uid, pwd, nil
}

// GetUser ...
func GetUser(uid int) (*defs.User, error) {
	stmt, err := dbConn.Prepare("select name from users where id = ?")
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("close db connection failed: %s", err)
		}
	}()
	if err != nil {
		log.Printf("(ERROR) GetUser sql prepare error: %s", err)
		return nil, err
	}
	var uname string
	err = stmt.QueryRow(uid).Scan(&uname)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("(ERROR) GetUser sql query error: %s", err)
		return nil, err
	}
	res := &defs.User{ID: uid, Name: uname}
	return res, nil
}

// DeleteUser ...
func DeleteUser(name string, pwd string) error {
	stmt, err := dbConn.Prepare("DELETE FROM users WHERE name = ? AND pwd = ?")
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("close db connection failed: %s", err)
		}
	}()
	if err != nil {
		log.Printf("(ERROR) DeleteUser sql prepare error: %s", err)
		return err
	}
	_, err = stmt.Exec(name, pwd)
	if err != nil {
		log.Printf("(ERROR) DeleteUser sql exec error: %s", err)
		return err
	}
	return nil
}

package dbops

import (
	"log"

	// Registe mysql
	_ "github.com/go-sql-driver/mysql"
)

// AddVideoDeletionRecord ...
func AddVideoDeletionRecord(vid int) error {
	stmt, err := dbConn.Prepare("INSERT INTO video_del_rec (video_id) VALUES(?)")
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("close db connection failed: %s", err)
		}
	}()
	if err != nil {
		return err
	}
	_, err = stmt.Exec(vid)
	return err
}

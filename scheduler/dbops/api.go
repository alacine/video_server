package dbops

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func AddVideoDeletionRecord(vid int) error {
	stmtIns, err := dbConn.Prepare("INSERT INTO video_del_rec (video_id) VALUES(?)")
	defer stmtIns.Close()
	if err != nil {
		return err
	}
	stmtIns.Exec(vid)
	if err != nil {
		log.Printf("AddVideoDeletionRecord error: %v", err)
		return err
	}
	return nil
}

package dbops

import (
	"database/sql"
	"log"
)

func ReadVideoDeletionRecord(count int) ([]int, error) {
	stmtOut, err := dbConn.Prepare("SELECT video_id FROM video_del_rec LIMIT ?")
	defer stmtOut.Close()
	var ids []int
	if err != nil {
		return ids, err
	}
	row, err := stmtOut.Query(count)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Query VideoDeletionRecord error: %v", err)
		return ids, err
	}
	for row.Next() {
		var id int
		if err := row.Scan(&id); err != nil {
			return ids, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func DelVideoDeletionRecord(vid int) error {
	stmtDel, err := dbConn.Prepare("DELETE FROM video_del_rec WHERE video_id = ?")
	defer stmtDel.Close()
	if err != nil {
		return err
	}
	_, err = stmtDel.Exec(vid)
	if err != nil {
		log.Printf("Deleting VideoDeletionRecord error: %v", err)
		return err
	}
	return nil
}

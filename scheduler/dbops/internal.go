package dbops

import (
	"database/sql"
	"log"
)

// ReadVideoDeletionRecord ...
func ReadVideoDeletionRecord(count int) ([]int, error) {
	stmt, err := dbConn.Prepare("SELECT video_id FROM video_del_rec LIMIT ?")
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("close db connection failed: %s", err)
		}
	}()
	var ids []int
	if err != nil {
		return ids, err
	}
	row, err := stmt.Query(count)
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

// DelVideoDeletionRecord ...
func DelVideoDeletionRecord(vid int) error {
	stmt, err := dbConn.Prepare("DELETE FROM video_del_rec WHERE video_id = ?")
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("close db connection failed: %s", err)
		}
	}()
	if err != nil {
		return err
	}
	_, err = stmt.Exec(vid)
	if err != nil {
		log.Printf("Deleting VideoDeletionRecord error: %v", err)
		return err
	}
	return nil
}

package dbops

import (
	"database/sql"
	"log"

	"github.com/alacine/video_server/api/defs"
	// Registe mysql
	_ "github.com/go-sql-driver/mysql"
)

// ListVideos ...
func ListVideos() ([]*defs.VideoInfo, error) {
	stmt, err := dbConn.Prepare(`
		SELECT
		  video_info.id, video_info.author_id, users.name,
		  video_info.title, video_info.create_time, video_info.description
		FROM video_info INNER JOIN users ON video_info.author_id = users.id
		ORDER BY video_info.create_time DESC
	`)
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("close db connection failed: %s", err)
		}
	}()
	var videos []*defs.VideoInfo
	if err != nil {
		log.Printf("(ERROR) ListVideos sql prepare error: %s", err)
		return videos, err
	}
	rows, err := stmt.Query()
	if err != nil {
		log.Printf("(ERROR) ListVideos sql query error: %s", err)
		return videos, err
	}
	for rows.Next() {
		var aname, title, ctime, desp string
		var id, aid int
		if err := rows.Scan(&id, &aid, &aname, &title, &ctime, &desp); err != nil {
			return videos, err
		}
		v := &defs.VideoInfo{
			ID:          id,
			AuthorID:    aid,
			AuthorName:  aname,
			Title:       title,
			CreateTime:  ctime,
			Description: desp,
		}
		videos = append(videos, v)
	}
	return videos, nil
}

// AddNewVideo ...
func AddNewVideo(aid int, title string, desp string) (*defs.VideoInfo, error) {
	if err != nil {
		return nil, err
	}
	stmt, err := dbConn.Prepare(`
		INSERT INTO video_info (author_id, title, description)
		VALUES(?, ?, ?)
	`)
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("close db connection failed: %s", err)
		}
	}()
	if err != nil {
		return nil, err
	}
	result, err := stmt.Exec(aid, title, desp)
	if err != nil {
		return nil, err
	}
	vid, err := result.LastInsertId()
	video := &defs.VideoInfo{ID: int(vid), AuthorID: aid, Title: title, Description: desp}
	return video, err
}

// GetVideoInfo ...
func GetVideoInfo(vid int) (*defs.VideoInfo, error) {
	stmt, _ := dbConn.Prepare(`
		SELECT
		  video_info.author_id, users.name, video_info.title,
		  video_info.create_time, video_info.description
		FROM video_info INNER JOIN users ON video_info.author_id = users.id
		WHERE video_info.id = ?
	`)
	var aid int
	var aname, ctime, title, desp string
	err = stmt.QueryRow(vid).Scan(&aid, &aname, &title, &ctime, &desp)
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("close db connection failed: %s", err)
		}
	}()
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	video := &defs.VideoInfo{
		ID:          vid,
		AuthorID:    aid,
		AuthorName:  aname,
		Title:       title,
		CreateTime:  ctime,
		Description: desp,
	}
	return video, nil
}

// ListUserVideos ...
func ListUserVideos(uid, from, to int) ([]*defs.VideoInfo, error) {
	stmt, err := dbConn.Prepare(`
		SELECT
		  video_info.id, video_info.author_id, users.name,
		  video_info.title, video_info.create_time, video_info.description
		FROM video_info INNER JOIN users ON video_info.author_id = users.id
		WHERE users.id = ?
		  AND video_info.create_time > FROM_UNIXTIME(?)
		  AND video_info.create_time <= FROM_UNIXTIME(?)
		ORDER BY video_info.create_time DESC
	`)
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("close db connection failed: %s", err)
		}
	}()
	var videos []*defs.VideoInfo
	if err != nil {
		log.Printf("(ERROR) ListUserVideos sql prepare error: %s", err)
		return videos, err
	}
	rows, err := stmt.Query(uid, from, to)
	if err != nil {
		log.Printf("(ERROR) ListUserVideos sql query error: %s", err)
		return videos, err
	}
	for rows.Next() {
		var title, aname, ctime, desp string
		var id, aid int
		if err := rows.Scan(&id, &aid, &aname, &title, &ctime, &desp); err != nil {
			return videos, err
		}
		v := &defs.VideoInfo{
			ID:          id,
			AuthorID:    aid,
			AuthorName:  aname,
			Title:       title,
			CreateTime:  ctime,
			Description: desp,
		}
		videos = append(videos, v)
	}
	return videos, nil
}

// DeleteVideoInfo ...
func DeleteVideoInfo(vid int) error {
	stmt, err := dbConn.Prepare("DELETE FROM video_info WHERE id = ?")
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("close db connection failed: %s", err)
		}
	}()
	if err != nil {
		log.Printf("(ERROR) DeleteVideoInfo sql prepare error: %s", err)
		return err
	}
	_, err = stmt.Exec(vid)
	if err != nil {
		log.Printf("(ERROR) DeleteVideoInfo sql exec error: %s", err)
		return err
	}
	return nil
}

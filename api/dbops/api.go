package dbops

import (
	"database/sql"
	"log"

	"github.com/alacine/video_server/api/defs"
	"github.com/alacine/video_server/api/utils"
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

func ListVideos() ([]*defs.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare(`
		SELECT
		  video_info.id, video_info.author_id, users.name,
		  video_info.title, video_info.create_time, video_info.description
		FROM video_info INNER JOIN users ON video_info.author_id = users.id
		ORDER BY video_info.create_time DESC
	`)
	defer stmtOut.Close()
	var videos []*defs.VideoInfo
	if err != nil {
		log.Printf("(ERROR) ListVideos sql prepare error: %s", err)
		return videos, err
	}
	rows, err := stmtOut.Query()
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
			Id:          id,
			AuthorId:    aid,
			AuthorName:  aname,
			Title:       title,
			CreateTime:  ctime,
			Description: desp,
		}
		videos = append(videos, v)
	}
	return videos, nil
}

func AddNewVideo(aid int, title string, desp string) (*defs.VideoInfo, error) {
	if err != nil {
		return nil, err
	}
	stmtIns, err := dbConn.Prepare(`
		INSERT INTO video_info (author_id, title, description) 
		VALUES(?, ?, ?)
	`)
	defer stmtIns.Close()
	if err != nil {
		return nil, err
	}
	result, err := stmtIns.Exec(aid, title, desp)
	if err != nil {
		return nil, err
	}
	vid, err := result.LastInsertId()
	video := &defs.VideoInfo{Id: int(vid), AuthorId: aid, Title: title, Description: desp}
	return video, nil
}

func GetVideoInfo(vid int) (*defs.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare(`
		SELECT
		  video_info.author_id, users.name, video_info.title,
		  video_info.create_time, video_info.description
		FROM video_info INNER JOIN users ON video_info.author_id = users.id
		WHERE video_info.id = ?
	`)
	var aid int
	var aname, ctime, title, desp string
	err = stmtOut.QueryRow(vid).Scan(&aid, &aname, &title, &ctime, &desp)
	defer stmtOut.Close()
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	video := &defs.VideoInfo{
		Id:          vid,
		AuthorId:    aid,
		AuthorName:  aname,
		Title:       title,
		CreateTime:  ctime,
		Description: desp,
	}
	return video, nil
}

func ListUserVideos(uid, from, to int) ([]*defs.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare(`
		SELECT
		  video_info.id, video_info.author_id, users.name,
		  video_info.title, video_info.create_time, video_info.description
		FROM video_info INNER JOIN users ON video_info.author_id = users.id
		WHERE users.id = ?
		  AND video_info.create_time > FROM_UNIXTIME(?)
		  AND video_info.create_time <= FROM_UNIXTIME(?)
		ORDER BY video_info.create_time DESC
	`)
	defer stmtOut.Close()
	var videos []*defs.VideoInfo
	if err != nil {
		log.Printf("(ERROR) ListUserVideos sql prepare error: %s", err)
		return videos, err
	}
	rows, err := stmtOut.Query(uid, from, to)
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
			Id:          id,
			AuthorId:    aid,
			AuthorName:  aname,
			Title:       title,
			CreateTime:  ctime,
			Description: desp,
		}
		videos = append(videos, v)
	}
	return videos, nil
}

func DeleteVideoInfo(vid int) error {
	stmtDel, err := dbConn.Prepare("DELETE FROM video_info WHERE id = ?")
	defer stmtDel.Close()
	if err != nil {
		log.Printf("(ERROR) DeleteVideoInfo sql prepare error: %s", err)
		return err
	}
	_, err = stmtDel.Exec(vid)
	if err != nil {
		log.Printf("(ERROR) DeleteVideoInfo sql exec error: %s", err)
		return err
	}
	return nil
}

func AddNewComment(vid int, aid int, content string) error {
	id, err := utils.NewUUID()
	if err != nil {
		return err
	}
	stmtIns, err := dbConn.Prepare(`
		INSERT INTO comments (id, video_id, author_id, content) 
		VALUES (?, ?, ?, ?)
	`)
	defer stmtIns.Close()
	if err != nil {
		log.Printf("(ERROR) AddNewComment sql prepare error: %s", err)
		return err
	}
	_, err = stmtIns.Exec(id, vid, aid, content)
	if err != nil {
		log.Printf("(ERROR) AddNewComment sql exec error: %s", err)
		return err
	}
	return nil
}

func ListComments(vid, from, to int) ([]*defs.Comment, error) {
	stmtOut, err := dbConn.Prepare(`
		SELECT comments.id, users.name, comments.content, comments.post_time
		FROM comments INNER JOIN users ON comments.author_id = users.id
		WHERE comments.video_id = ?
		  AND comments.post_time > FROM_UNIXTIME(?)
		  AND comments.post_time <= FROM_UNIXTIME(?)
		ORDER BY comments.post_time DESC
	`)
	/* 注意这里查询的区间是前开后闭，后带等号是因为在 MYSQL 里面记录的时间到秒，
	 * 如果 to 是当前时间而且是开区间，写入之后马上读取会发生读不到的情况
	 */
	defer stmtOut.Close()
	var comments []*defs.Comment
	if err != nil {
		log.Printf("(ERROR) ListComments sql prepare error: %s", err)
		return comments, err
	}
	rows, err := stmtOut.Query(vid, from, to)
	if err != nil {
		log.Printf("(ERROR) ListComments sql exec error: %s", err)
		return comments, err
	}
	for rows.Next() {
		var id, name, content, ptime string
		if err := rows.Scan(&id, &name, &content, &ptime); err != nil {
			return comments, err
		}
		c := &defs.Comment{
			Id:         id,
			VideoId:    vid,
			AuthorName: name,
			Content:    content,
			PostTime:   ptime,
		}
		comments = append(comments, c)
	}
	return comments, nil
}

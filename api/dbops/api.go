package dbops

import (
	"database/sql"
	"log"
	"time"

	"github.com/alacine/video_server/api/defs"
	"github.com/alacine/video_server/api/utils"
	_ "github.com/go-sql-driver/mysql"
)

func AddUserCredential(loginName string, pwd string) error {
	stmtIns, err := dbConn.Prepare("INSERT INTO users (login_name, pwd) VALUES (?, ?)")
	defer stmtIns.Close()
	if err != nil {
		log.Printf("AddUserCredential error %s", err)
		return err
	}
	_, err = stmtIns.Exec(loginName, pwd)
	if err != nil {
		log.Printf("AddUserCredential error %s", err)
		return err
	}
	return nil
}

func GetUserCredential(loginName string) (string, error) {
	stmtOut, err := dbConn.Prepare("SELECT pwd FROM users WHERE login_name = ?")
	defer stmtOut.Close()
	if err != nil {
		log.Printf("GetUserCredential error %s", err)
		return "", err
	}
	var pwd string
	err = stmtOut.QueryRow(loginName).Scan(&pwd)
	if err != nil && nil != sql.ErrNoRows {
		return "", err
	}
	return pwd, nil
}

func DeleteUser(loginName string, pwd string) error {
	stmtDel, err := dbConn.Prepare("DELETE FROM users WHERE login_name = ? AND pwd = ?")
	defer stmtDel.Close()
	if err != nil {
		log.Printf("DeleteUser error: %s", err)
		return err
	}
	_, err = stmtDel.Exec(loginName, pwd)
	if err != nil {
		return err
	}
	return nil
}

func AddNewVideo(aid int, name string) (*defs.VideoInfo, error) {
	// create uuid
	vid, err := utils.NewUUID()
	if err != nil {
		return nil, err
	}
	t := time.Now()
	ctime := t.Format("Jan 02 2006, 15:04:05") //M D y, HH:MM:SS
	stmtIns, err := dbConn.Prepare(`INSERT INTO video_info (id, author_id, name, display_ctime) 
									VALUES(?, ?, ?, ?)`)
	defer stmtIns.Close()
	if err != nil {
		return nil, err
	}
	_, err = stmtIns.Exec(vid, aid, name, ctime)
	if err != nil {
		return nil, err
	}
	video := &defs.VideoInfo{Id: vid, AuthorId: aid, Name: name, DisplayCtime: ctime}
	return video, nil
}

func GetVideoInfo(vid string) (*defs.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare(`SELECT author_id, name, display_ctime 
									FROM video_info WHERE id = ?`)
	var aid int
	var dct string
	var name string
	err = stmtOut.QueryRow(vid).Scan(&aid, &name, &dct)
	defer stmtOut.Close()
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	video := &defs.VideoInfo{Id: vid, AuthorId: aid, Name: name, DisplayCtime: dct}
	return video, nil
}

func DeleteVideoInfo(vid string) error {
	stmtDel, err := dbConn.Prepare("DELETE FROM video_info WHERE id = ?")
	defer stmtDel.Close()
	if err != nil {
		return err
	}
	_, err = stmtDel.Exec(vid)
	if err != nil {
		return err
	}
	return nil
}

func AddNewComment(vid string, aid int, content string) error {
	id, err := utils.NewUUID()
	if err != nil {
		return err
	}
	stmtIns, err := dbConn.Prepare(`INSERT INTO comments (id, video_id, author_id, content) 
									VALUES (?, ?, ?, ?)`)
	defer stmtIns.Close()
	if err != nil {
		return err
	}
	_, err = stmtIns.Exec(id, vid, aid, content)
	if err != nil {
		return err
	}
	return nil
}

func ListComments(vid string, from, to int) ([]*defs.Comment, error) {
	stmtOut, err := dbConn.Prepare(`SELECT comments.id, users.login_name, comments.content
									FROM comments INNER JOIN users ON comments.author_id = users.id
									WHERE comments.video_id = ?
									  AND comments.time > FROM_UNIXTIME(?)
									  AND comments.time <= FROM_UNIXTIME(?)`)
	/* 注意这里查询的区间是前开后闭，后带等号是因为在 MYSQL 里面记录的时间到秒，
	 * 如果 to 是当前时间而且是开区间，写入之后马上读取会发生读不到的情况
	 */
	defer stmtOut.Close()
	var comments []*defs.Comment
	rows, err := stmtOut.Query(vid, from, to)
	if err != nil {
		return comments, err
	}
	for rows.Next() {
		var id, name, content string
		if err := rows.Scan(&id, &name, &content); err != nil {
			return comments, err
		}
		c := &defs.Comment{Id: id, VideoId: vid, AuthorName: name, Content: content}
		comments = append(comments, c)
	}
	return comments, nil
}

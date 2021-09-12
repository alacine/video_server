package dbops

import (
	"log"

	"github.com/alacine/video_server/api/defs"
	"github.com/alacine/video_server/api/utils"

	// Registe mysql
	_ "github.com/go-sql-driver/mysql"
)

// AddNewComment ...
func AddNewComment(vid int, aid int, content string) error {
	id, err := utils.NewUUID()
	if err != nil {
		return err
	}
	stmt, err := dbConn.Prepare(`
		INSERT INTO comments (id, video_id, author_id, content)
		VALUES (?, ?, ?, ?)
	`)
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("close db connection failed: %s", err)
		}
	}()
	if err != nil {
		log.Printf("(ERROR) AddNewComment sql prepare error: %s", err)
		return err
	}
	_, err = stmt.Exec(id, vid, aid, content)
	if err != nil {
		log.Printf("(ERROR) AddNewComment sql exec error: %s", err)
		return err
	}
	return nil
}

// ListComments ...
func ListComments(vid, from, to int) ([]*defs.Comment, error) {
	stmt, err := dbConn.Prepare(`
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
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("close db connection failed: %s", err)
		}
	}()
	var comments []*defs.Comment
	if err != nil {
		log.Printf("(ERROR) ListComments sql prepare error: %s", err)
		return comments, err
	}
	rows, err := stmt.Query(vid, from, to)
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
			ID:         id,
			VideoID:    vid,
			AuthorName: name,
			Content:    content,
			PostTime:   ptime,
		}
		comments = append(comments, c)
	}
	return comments, nil
}

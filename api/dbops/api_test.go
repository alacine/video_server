package dbops

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
	"time"
)

var tmpvid int

func clearTables() {
	_, err := dbConn.Exec("truncate users")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = dbConn.Exec("truncate video_info")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = dbConn.Exec("truncate comments")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = dbConn.Exec("truncate sessions")
	if err != nil {
		log.Fatalln(err)
	}
}

func TestMain(m *testing.M) {
	//clearTables()
	//m.Run()
	//clearTables()
	os.Exit(0)
}

func TestUserWorkFlow(t *testing.T) {
	t.Run("Add", testAddUser)
	t.Run("Get", testGetUser)
	t.Run("Del", testDeleteUser)
	t.Run("Reget", testRegetUser)
}

func testAddUser(t *testing.T) {
	_, err := AddUserCredential("test_man", "abc")
	if err != nil {
		t.Errorf("Error of AddUser: %v", err)
	}
}

func testGetUser(t *testing.T) {
	_, pwd, err := GetUserCredential("test_man")
	if pwd != "abc" || err != nil {
		t.Errorf("Error of GetUser")
	}
}

func testDeleteUser(t *testing.T) {
	err := DeleteUser("test_man", "abc")
	if err != nil {
		t.Errorf("Error of DeleteUser: %v", err)
	}
}

func testRegetUser(t *testing.T) {
	_, pwd, err := GetUserCredential("test_man")
	if err != nil && err != sql.ErrNoRows {
		t.Errorf("Error of RegetUser: %v", err)
	}

	if pwd != "" {
		t.Errorf("Deleting user test failed")
	}
}

func TestVideoWorkFlow(t *testing.T) {
	clearTables()
	t.Run("PrepareUser", testAddUser)
	t.Run("AddVideo", testAddVideoInfo)
	t.Run("GetVideo", testGetVideoInfo)
	t.Run("DeleteVideo", testDeleteVideoInfo)
	t.Run("RegetVideo", testRegetVideoInfo)
}

func testAddVideoInfo(t *testing.T) {
	video, err := AddNewVideo(1, "test-video-first", "description-video-first")
	if err != nil {
		t.Errorf("Error of AddVideoInfo: %v", err)
	}
	tmpvid = video.ID
	tmpvid = video.ID
}

func testGetVideoInfo(t *testing.T) {
	video, err := GetVideoInfo(tmpvid)
	if video.Title != "test-video-first" || err != nil {
		t.Errorf("Error of GetVideoInfo: %v", err)
	}
}

func testDeleteVideoInfo(t *testing.T) {
	err := DeleteVideoInfo(tmpvid)
	if err != nil {
		t.Errorf("Error of DeleteVideoInfo: %v", err)
	}
}

func testRegetVideoInfo(t *testing.T) {
	video, err := GetVideoInfo(tmpvid)
	if video != nil && err != sql.ErrNoRows {
		t.Errorf("Error of RegetVideoInfo: %v", err)
	}
}

func TestCommets(t *testing.T) {
	clearTables()
	t.Run("AddUser", testAddUser)
	t.Run("AddComment", testAddComment)
	t.Run("ListComments", testListComments)
}

func testAddComment(t *testing.T) {
	vid := 123
	aid := 1
	content := "I like this video"
	err := AddNewComment(vid, aid, content)
	if err != nil {
		t.Errorf("Error of AddComment: %v", err)
	}
}

func testListComments(t *testing.T) {
	vid := 123
	from := 1514764800
	to, _ := strconv.Atoi(strconv.FormatInt(time.Now().UnixNano()/1e9, 10))
	comments, err := ListComments(vid, from, to)
	if err != nil {
		t.Errorf("Error of ListComments: %v", err)
	}
	for i, c := range comments {
		fmt.Printf("comment: %d, %v\n", i, c)
	}
}

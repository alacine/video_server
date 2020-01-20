package dbops

import (
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"
)

var tmpvid string

func clearTables() {
	dbConn.Exec("truncate users")
	dbConn.Exec("truncate video_info")
	dbConn.Exec("truncate comments")
	dbConn.Exec("truncate sessions")
}

func TestMain(m *testing.M) {
	clearTables()
	m.Run()
	clearTables()
}

func TestUserWorkFlow(t *testing.T) {
	t.Run("Add", testAddUser)
	t.Run("Get", testGetUser)
	t.Run("Del", testDeleteUser)
	t.Run("Reget", testRegetUser)
}

func testAddUser(t *testing.T) {
	err := AddUserCredential("ryan", "abc")
	if err != nil {
		t.Errorf("Error of AddUser: %v", err)
	}
}

func testGetUser(t *testing.T) {
	pwd, err := GetUserCredential("ryan")
	if pwd != "abc" || err != nil {
		t.Errorf("Error of GetUser")
	}
}

func testDeleteUser(t *testing.T) {
	err := DeleteUser("ryan", "abc")
	if err != nil {
		t.Errorf("Error of DeleteUser: %v", err)
	}
}

func testRegetUser(t *testing.T) {
	pwd, err := GetUserCredential("ryan")
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
	video, err := AddNewVideo(1, "test-video-first")
	if err != nil {
		t.Errorf("Error of AddVideoInfo: %v", err)
	}
	tmpvid = video.Id
}

func testGetVideoInfo(t *testing.T) {
	video, err := GetVideoInfo(tmpvid)
	if video.Name != "test-video-first" || err != nil {
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
	vid := "123"
	aid := 1
	content := "I like this video"
	err := AddNewComment(vid, aid, content)
	if err != nil {
		t.Errorf("Error of AddComment: %v", err)
	}
}

func testListComments(t *testing.T) {
	vid := "123"
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

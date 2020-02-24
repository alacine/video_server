package defs

// requests
type UserCredential struct {
	Username string `json:"user_name"`
	Pwd      string `json:"pwd"`
}

type NewComment struct {
	AuthorId int    `json:"author_id"`
	Content  string `json:"content"`
}

type NewVideo struct {
	AuthorId int    `json:"author_id"`
	Name     string `json:"name"`
}

// response
type SignedUp struct {
	Success   bool   `json:"success"`
	SessionId string `json:"session_id"`
}

type SignedIn struct {
	Success   bool   `json:"success"`
	SessionId string `json:"session_id"`
}

type UserInfo struct {
	Username string `json:"user_name"`
	Id       int    `json:"user_id"`
}

type VideosInfo struct {
	Videos []*VideoInfo `json:"videos"`
}

type Comments struct {
	Comment []*Comment `json:"comments"`
}

// Data model
type User struct {
	Id        int    `json:"user_id"`
	LoginName string `json:"login_name"`
	Pwd       string `json:"pwd"`
}

type VideoInfo struct {
	Id           int
	AuthorId     int
	Name         string
	DisplayCtime string
	Description  string `json:"description"`
}

type Comment struct {
	Id         string
	VideoId    int
	AuthorName string
	Content    string
}

type SimpleSession struct {
	Username string // login name
	TTL      int64
}

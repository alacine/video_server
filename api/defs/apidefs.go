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
	AuthorId    int    `json:"author_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Page struct {
	Start int `json:"start"`
	Limit int `json:"Page"`
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
	Id   int    `json:"user_id"`
	Name string `json:"name"`
	Pwd  string `json:"pwd"`
}

type VideoInfo struct {
	Id           int    `json:"id"`
	AuthorId     int    `json:"author_id"`
	AuthorName   string `json:"author_name"`
	Title        string `json:"title"`
	DisplayCtime string `json:"displayctime"`
	Description  string `json:"description"`
}

type Comment struct {
	Id         string `json:"id"`
	VideoId    int    `json:"video_id"`
	AuthorName string `json:"author_name"`
	Content    string `json:"content"`
}

type SimpleSession struct {
	Username string // login name
	TTL      int64
}

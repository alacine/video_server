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
//type SignedUp struct {
//Success   bool   `json:"success"`
//SessionId string `json:"session_id"`
//}

type SignedIn struct {
	Success   bool   `json:"success"`
	SessionId string `json:"session_id"`
	UserId    int    `json:"user_id"`
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
	Id          int    `json:"id"`
	AuthorId    int    `json:"author_id"`
	AuthorName  string `json:"author_name"`
	Title       string `json:"title"`
	CreateTime  string `json:"create_time"`
	Description string `json:"description"`
}

type Comment struct {
	Id         string `json:"id"`
	VideoId    int    `json:"video_id"`
	AuthorName string `json:"author_name"`
	Content    string `json:"content"`
	PostTime   string `json:"post_time"`
}

type SimpleSession struct {
	SessionId string
	UserId    int
	TTL       int64
}

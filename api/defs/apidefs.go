package defs

// UserCredential requests
type UserCredential struct {
	Username string `json:"user_name"`
	Pwd      string `json:"pwd"`
}

// NewComment ...
type NewComment struct {
	AuthorID int    `json:"author_id"`
	Content  string `json:"content"`
}

// NewVideo ...
type NewVideo struct {
	AuthorID    int    `json:"author_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Page ...
type Page struct {
	Start int `json:"start"`
	Limit int `json:"Page"`
}

// response
//type SignedUp struct {
//Success   bool   `json:"success"`
//SessionID string `json:"session_id"`
//SessionID string `json:"session_id"`
//}

// SignedIn ...
type SignedIn struct {
	Success   bool   `json:"success"`
	SessionID string `json:"session_id"`
	UserID    int    `json:"user_id"`
}

// UserInfo ...
type UserInfo struct {
	Username string `json:"user_name"`
	ID       int    `json:"user_id"`
}

// VideosInfo ...
type VideosInfo struct {
	Videos []*VideoInfo `json:"videos"`
}

// Comments ...
type Comments struct {
	Comment []*Comment `json:"comments"`
}

// User ...
type User struct {
	ID   int    `json:"user_id"`
	Name string `json:"name"`
	Pwd  string `json:"pwd"`
}

// VideoInfo ...
type VideoInfo struct {
	ID          int    `json:"id"`
	AuthorID    int    `json:"author_id"`
	AuthorName  string `json:"author_name"`
	Title       string `json:"title"`
	CreateTime  string `json:"create_time"`
	Description string `json:"description"`
}

// Comment ...
type Comment struct {
	ID         string `json:"id"`
	VideoID    int    `json:"video_id"`
	AuthorName string `json:"author_name"`
	Content    string `json:"content"`
	PostTime   string `json:"post_time"`
}

// SimpleSession ...
type SimpleSession struct {
	SessionID string
	UserID    int
	TTL       int64
}

package domain

import (
	"time"
)

type User struct {
	Nickname string `json:"nickname"`
	About    string `json:"about"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
}

type Forum struct {
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Owner   string `json:"user"`
	Threads int32  `json:"threads"`
	Posts   int32  `json:"posts"`
}

type Thread struct {
	Id      int       `json:"id"`
	Author  string    `json:"author"`
	Created time.Time `json:"created"`
	Forum   string    `json:"forum"`
	Message string    `json:"message"`
	Slug    string    `json:"slug"`
	Title   string    `json:"title"`
	Votes   int       `json:"votes"`
}

type Post struct {
	Id        int       `json:"id"`
	Author    string    `json:"author"`
	Created   time.Time `json:"created"`
	ForumSlug string    `json:"forum"`
	IsEdited  bool      `json:"isEdited"`
	Message   string    `json:"message"`
	Parent    int       `json:"parent,omitempty"`
	Posts     []int     `json:"-"`
	Thread    int       `json:"thread"`
}

type PostFull struct {
	Profile *User   `json:"author,omitempty"`
	Forum   *Forum  `json:"forum,omitempty"`
	Post    Post    `json:"post"`
	Thread  *Thread `json:"thread,omitempty"`
}

type Vote struct {
	Profile         int    `json:"-"`
	ProfileNickname string `json:"nickname"`
	Thread          int    `json:"-"`
	Voice           int    `json:"voice"`
}

type Status struct {
	Forum  int `json:"forum"`
	Post   int `json:"post"`
	Thread int `json:"thread"`
	User   int `json:"user"`
}
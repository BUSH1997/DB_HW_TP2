package models

import "time"

type Forum struct {
	Title   string `json:"title" validate:"required"`
	User    string `json:"user" validate:"required"`
	Slug    string `json:"slug" validate:"required"`
	Posts   int64  `json:"posts"`
	Threads int64  `json:"threads"`
}

type Thread struct {
	Id      int32     `json:"id"`
	Title   string    `json:"title" validate:"required"`
	Author  string    `json:"author" validate:"required"`
	Forum   string    `json:"forum"`
	Message string    `json:"message" validate:"required"`
	Votes   int32     `json:"votes"`
	Slug    string    `json:"slug,omitempty"`
	Created time.Time `json:"created"`
}

type Vote struct {
	NickName string `json:"nickname"`
	Voice    int    `json:"voice"`
}

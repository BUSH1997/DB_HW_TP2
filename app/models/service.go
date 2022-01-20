package models

type Status struct {
	User int64 `json:"user,omitempty"`
	Forum int64 `json:"forum,omitempty"`
	Thread int64 `json:"thread,omitempty"`
	Post int64 `json:"post,omitempty"`
}

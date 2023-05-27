package database

import "net/url"

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Link struct {
	ID   int64   `json:"id"`
	Name string  `json:"name"`
	URL  url.URL `json:"link"`
	User User    `json:"user"`
}

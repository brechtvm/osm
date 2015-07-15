package user

import ()

type User struct {
	Id   int64
	Name string
}

func New(id int64, name string) *User {
	return &User{Id: id, Name: name}
}

// vim: ts=4 sw=4 noexpandtab nolist syn=go

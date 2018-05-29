package user

import ()

type User struct {
	Id   uint32
	Name string
}

func New(id uint32, name string) *User {
	return &User{Id: id, Name: name}
}

// vim: ts=4 sw=4 noexpandtab nolist syn=go

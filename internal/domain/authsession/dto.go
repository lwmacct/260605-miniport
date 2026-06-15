package authsession

import "time"

type Request struct {
	IP         string
	Host       string
	UserAgent  string
	Method     string
	Path       string
	RemoteAddr string
}

type UserDTO struct {
	ID        int64
	Username  string
	ExpiresAt time.Time
}

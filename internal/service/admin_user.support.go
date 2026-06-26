package service

import "time"

type AdminUser struct {
	ID          int64
	Username    string
	DisplayName string
	Status      string
	Admin       bool
	DisabledAt  *time.Time
}

func utilAdminUser(user IdentityUser, admin bool) AdminUser {
	return AdminUser{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Status:      user.Status,
		Admin:       admin,
		DisabledAt:  user.DisabledAt,
	}
}

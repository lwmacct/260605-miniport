package handler

import (
	"context"

	"github.com/lwmacct/260605-miniport/internal/service"
)

type Config struct {
	LocalLoginEnabled        bool
	LocalRegistrationEnabled bool
	SecureCookies            bool
	RuntimeAdmins            []string
	Request                  RequestAuth
}

type Services struct {
	Users      *service.UserService
	Passwords  *service.AuthPasswordService
	Sessions   *service.AuthSessionService
	Challenges *service.AuthChallengeService
	AdminUsers *service.AdminUserService
	Portsvc    *service.PortsvcService
}

type RequestAuth func(context.Context) (service.AuthSessionInput, bool)

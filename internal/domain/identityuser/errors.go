package identityuser

import "errors"

var (
	ErrDisabled           = errors.New("user disabled")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUsernameTaken      = errors.New("username taken")
)

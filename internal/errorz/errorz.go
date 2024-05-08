package errorz

import (
	"errors"
)

var (
	TokenExpired = errors.New("tokenExpiredError")
	UserExists   = errors.New("userAlreadyRegisteredError")
)

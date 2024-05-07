package errorz

import (
	"errors"
)

var (
	TokenExpired = errors.New("token expired")
)

package exceptions

import "errors"

var (
	ErrUserNotLoggedIn = errors.New("ErrUserNotLoggedIn")
)

type ErrLoginError struct {
	Reason string
}

func (ErrLoginError) Error() string {
	return "ErrLoginError"
}

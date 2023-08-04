package token

type ErrorType string

const (
	Expired ErrorType = "token has expired"
	Invalid ErrorType = "token is invalid"
)

type Error struct {
	t ErrorType
}

var (
	ExpiredTokenErr = &Error{t: Expired}
	InvalidTokenErr = &Error{t: Invalid}
)

func (e *Error) Error() string {
	return string(e.t)
}

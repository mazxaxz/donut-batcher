package rest

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NewError(code string, err error) Error {
	e := Error{
		Code:    code,
		Message: err.Error(),
	}
	return e
}

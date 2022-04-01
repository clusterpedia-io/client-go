package utils

type E struct {
	s string
}

func (e *E) Error() string {
	return e.s
}

func (e *E) SetError(s string) {
	e.s = s
}

func NewError(s string) *E {
	return &E{s: s}
}

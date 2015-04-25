package ehttp

import "fmt"

type (
	MissingParamError string

	BadParamError struct {
		Param  string
		Reason string
	}
)

func NewMissingParamError(p string) *MissingParamError {
	e := MissingParamError(p)
	return &e
}

func (m *MissingParamError) Error() string {
	return fmt.Sprintf("handles error: missing param %s", string(*m))
}

func NewBadParamError(p string, r string) *BadParamError {
	return &BadParamError{
		Param:  p,
		Reason: r,
	}
}

func (b *BadParamError) Error() string {
	return fmt.Sprintf("handles error: bad param %s. %s", b.Param, b.Reason)
}

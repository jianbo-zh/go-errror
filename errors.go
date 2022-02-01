package errors

import (
	"bytes"
	"errors"
)

type Context = map[string]interface{}

type Error struct {
	Inner error
	Cod   int
	Msg   string
	Ctx   Context
}

func (err *Error) Error() string {
	if err.Inner == nil {
		return err.Msg
	}

	buff := bufferPool.Get().(*bytes.Buffer)
	buff.Reset()

	buff.WriteString(err.Msg)
	buff.Write(innerSeparator)
	buff.WriteString(err.Inner.Error())

	result := buff.String()
	bufferPool.Put(buff)

	return result
}

func (err *Error) Unwrap() error {
	return err.Inner
}

func (err *Error) Code() int {
	return err.Cod
}

func (err *Error) Data() Context {
	return err.Ctx
}

func New(msg string, options ...Option) error {

	err := &Error{
		Msg: msg,
	}

	for _, option := range options {
		option(err)
	}

	return err
}

func NewC(code int, msg string, options ...Option) error {

	err := &Error{
		Cod: code,
		Msg: msg,
	}

	for _, option := range options {
		option(err)
	}

	return err
}

func Data(err error) Context {
	u, ok := err.(interface{ Data() Context })
	if !ok {
		return nil
	}
	return u.Data()
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

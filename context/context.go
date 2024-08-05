package context

import (
	"errors"

	"github.com/labstack/echo/v4"
)

var (
	ErrBindContextFailed         = errors.New("binding context failed")
	ErrRequestModelTypeAssertion = errors.New("request model type assertion failed")
)

type Context struct {
	echo.Context
	RequestModel any
}

func GetRequestModel[T any](ctx *Context) (T, error) {
	if ctx.RequestModel != nil {
		if value, ok := ctx.RequestModel.(T); ok {
			return value, nil
		} else {
			return value, ErrRequestModelTypeAssertion
		}
	} else {
		var zv T
		return zv, errors.New("type assertion failed")
	}
}

func Binder(echoContext echo.Context, i any) (*Context, error) {
	if ctx, ok := echoContext.(*Context); !ok {
		return nil, ErrBindContextFailed
	} else if err := ctx.Bind(i); err != nil {
		return nil, err
	} else {
		ctx.RequestModel = i
		return ctx, nil
	}
}

package apiresp

import (
	"backend/internal/api/apiresp/errs"
	"errors"
)

type ApiResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

func ApiSuccess(data any) *ApiResponse {
	return &ApiResponse{Data: data}
}

func ParseError(err error) *ApiResponse {
	if err == nil {
		return ApiSuccess(nil)
	}
	var codeErr *errs.CodeError
	if !errors.As(err, &codeErr) {
		codeErr = errs.ErrInternalServer.WithDetail(err.Error())
	}
	return &ApiResponse{Code: codeErr.Code, Msg: codeErr.Msg}
}

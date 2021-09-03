package controller

import (
	errors "github.com/xunlbz/go-restful-template/pkg/errors"
)

//all errors here
var (
	errorNoAuth           = errors.New(1001, "token无效")
	errorExpiredAuth      = errors.New(1002, "token失效")
	errorParseParams      = errors.New(1003, "参数解析错误")
	errorSetParams        = errors.New(1004, "参数错误")
	errorNoModule         = errors.New(1005, "指定模块不存在")
	errorDockerRestart    = errors.New(1006, "容器重启失败")
	errorServiceRestart   = errors.New(1007, "服务重启失败")
	errorSystemError      = errors.New(1008, "系统异常")
	errorUsernamePassword = errors.New(1009, "用户名或密码无效")
	errorUserExists       = errors.New(1010, "管理员已存在")
	errorUserNotExists    = errors.New(1011, "用户名不存在")
)

func NewSystemError(err error) error {
	return errors.New(1008, err.Error())
}

package handler

import (
	"context"
	"encoding/json"
	"github.com/micro/go-micro/util/log"
	protoUserApi "github.com/octopuszy/go-micro-api/proto/userApi"
	"github.com/octopuszy/go-micro-user/proto/user"
)

type UserApi struct {
	UserService user.UserService
}

func (u *UserApi) Login(ctx context.Context, req *protoUserApi.Request, res *protoUserApi.Response) error{
	log.Info("接收到请求...")
	username	:= req.Get["UserName"].Values[0]
	password	:= req.Get["Password"].Values[0]

	// 调用user服务，注册用户
	login := user.LoginReq{UserName: username, Password: password}
	rsp, err := u.UserService.Login(ctx, &login)
	if err != nil {
		res.Code = 1
		return err
	}
	repJson, _ := json.Marshal(rsp)
	res.Body = string(repJson)
	res.Code = 0
	return nil
}

func (u *UserApi) Register(ctx context.Context, req *protoUserApi.Request, res *protoUserApi.Response) error{
	log.Info("接收到请求...")
	username	:= req.Get["UserName"].Values[0]
	email 		:= req.Get["Email"].Values[0]
	password	:= req.Get["Password"].Values[0]

	// 调用user服务，注册用户
	register := user.RegisterReq{UserName: username, Email: email, Password: password}
	rsp, err := u.UserService.Register(ctx, &register)
	if err != nil {
		res.Code = 1
		return err
	}
	repJson, _ := json.Marshal(rsp)
	res.Body = string(repJson)
	res.Code = 0
	return nil
}
package users

import (
	"context"
)

type Service interface {
	CreateUser(context.Context, *ReqCreateUser) error
	GetUserInfo(context.Context, *ReqGetUserInfo) (*User, error)
	Login(context.Context, *ReqLogin) (*Token, error)
	Logout(context.Context, *ReqLogout) error
	ValidateToken(context.Context, *ReqValidateToken) (*TokenMeta, error)
	UpdateUser(context.Context, *ReqUpdateUserInfo) error
}

package users

import (
	"fmt"

	"github.com/Clay294/forum/exception"
	"github.com/Clay294/forum/flog"
	"github.com/golang-jwt/jwt/v5"
)

type ReqCreateUser struct {
	*ReqCreateUserBase
	*ReqCreateUserMeta
}

type ReqCreateUserBase struct {
	Name     string `json:"name,omitempty" validate:"required,min=1,max=25,validUserInfo"`
	Mail     string `json:"mail,omitempty" validate:"required,validUserInfo"`
	Password string `json:"password,omitempty" validate:"required,min=1,max=25,validUserInfo"`
}

type ReqCreateUserMeta struct {
	CreatedIn string `json:"created_in"`
}

type User struct {
	*UserBase
	*UserMeta
}

type UserBase ReqCreateUserBase

type UserMeta struct {
	Id             int    `json:"id,omitempty"`
	Coin           int    `json:"coin,omitempty"`
	CreatedAt      int64  `json:"created_at,omitempty"`
	CreatedIn      string `json:"created_in,omitempty"`
	UpdatedAt      int64  `json:"updated_at,omitempty"`
	UpdatedIn      string `json:"updated_in,omitempty"`
	LastLoggedinAt int64  `json:"last_loggedin_at,omitempty"`
	LastLoggedinIn string `json:"last_loggedin_in,omitempty"`
}

type ReqGetUserInfo struct {
	UserId int `json:"user_id" validate:"required,number"`
}

type ReqLogin struct {
	Credential string `json:"credential" validate:"required"`
	Password   string `json:"password" validate:"required"`
	*ReqLoginMeta
}

type ReqLogout struct {
	JsonTokenId string `json:"jti"`
}

type ReqLoginMeta struct {
	LoginBy    LOGINBY
	LoggedinIn string
	LoggedinAt int64
}

type Token struct {
	JsonToken string `json:"json_token,omitempty"`
	*TokenMeta
}

type TokenMeta struct {
	Id         int    `json:"id"`
	Issuer     string `json:"iss,omitempty"`
	Subject    int    `json:"sub,omitempty"`
	UserName   string `json:"user_name,omitempty"`
	UserMail   string `json:"user_mail,omitempty"`
	ExpiresAt  int64  `json:"exp,omitempty"`
	NotBefore  int64  `json:"nbf,omitempty"`
	IssuedAt   int64  `json:"iat,omitempty"`
	ID         string `json:"jti,omitempty" gorm:"column:uuid"`
	LoggedinIn string `json:"loggedin_in,omitempty"`
}

type JsonTokenClaims struct {
	*jwt.RegisteredClaims
	UserName   string `json:"user_name,omitempty"`
	UserMail   string `json:"user_mail,omitempty"`
	LoggedinIn string `json:"loggedin_in,omitempty"`
}

type ReqValidateToken struct {
	JsonToken string `json:"json_token" valdiate:"jwt"`
	*ReqValidateTokenMeta
}

type ReqValidateTokenMeta struct {
	ClientIP string `json:"client_ip"`
}

type TokenInBlacklist struct {
	Id   int    `json:"id"`
	UUID string `json:"uuid"`
}

type ReqUpdateUserInfo struct {
	Name     string `json:"name" validate:"required_without_all=Mail Password,excluded_with=Mail Password,omitempty,min=1,max=25,validUserInfo"`
	Mail     string `json:"mail" validate:"required_without_all=Name Password,excluded_with=Name Password,omitempty,validUserInfo"`
	Password string `json:"password" validate:"required_without_all=Name Mail,excluded_with=Name Mail,omitempty,min=1,max=25,validUserInfo"`
	*ReqUpdateUserInfoMeta
}

type ReqUpdateUserInfoMeta struct {
	UserId      int    `json:"user_id"`
	UpdatedAt   int64  `json:"updated_at"`
	UpdatedIn   string `json:"updated_in"`
	JsonTokenId string `json:"json_token_id"`
}

func (*User) TableName() string {
	return UnitName
}

func (*Token) TableName() string {
	return TokensTableName
}

func (*TokenMeta) TableName() string {
	return TokensTableName
}

func (rl *ReqLogin) SetLoginBy() error {
	isMail, err := mRE.MatchString(rl.Credential)

	if err != nil {
		flog.Flogger().Error().Msgf("unit %s, verifying credential from login request timeout: %s", UnitName, err)
		return exception.NewErrLogin(fmt.Errorf("登陆失败"))
	}

	if isMail {
		rl.LoginBy = LoginByMail
		return nil
	}

	rl.LoginBy = LoginByName
	return nil
}

func (*TokenInBlacklist) TableName() string {
	return TokensBlacklistTableName
}

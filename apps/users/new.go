package users

import (
	"strconv"
	"time"

	"github.com/Clay294/forum/protocol"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func NewReqCreateUser() *ReqCreateUser {
	return &ReqCreateUser{
		// ReqCreateUserBase: new(ReqCreateUserBase),
		ReqCreateUserBase: &ReqCreateUserBase{},
		// ReqCreateUserMeta: new(ReqCreateUserMeta),
		ReqCreateUserMeta: &ReqCreateUserMeta{},
	}
}

func NewUser(rcu *ReqCreateUser) *User {
	return &User{
		UserBase: (*UserBase)(rcu.ReqCreateUserBase),
		UserMeta: NewUserMeta(rcu.ReqCreateUserMeta),
	}
}

func NewUserMeta(rcum *ReqCreateUserMeta) *UserMeta {
	now := time.Now().UnixMilli()
	return &UserMeta{
		Coin:      500,
		CreatedAt: now,
		UpdatedAt: now,
		CreatedIn: rcum.CreatedIn,
		UpdatedIn: rcum.CreatedIn,
	}
}

func NewReqLogin() *ReqLogin {
	return &ReqLogin{
		ReqLoginMeta: NewReqLoginMeta(),
	}
}

func NewReqLoginMeta() *ReqLoginMeta {
	return &ReqLoginMeta{
		LoggedinAt: time.Now().UnixMilli(),
	}
}

func NewToken(user *User) *Token {
	return &Token{
		TokenMeta: NewTokenMeta(user),
	}
}

func NewTokenMeta(user *User) *TokenMeta {
	now := time.Now()
	return &TokenMeta{
		Issuer:     protocol.Domain,
		Subject:    user.Id,
		UserName:   user.Name,
		UserMail:   user.Mail,
		ExpiresAt:  now.Add(time.Hour * 24).UnixMilli(),
		NotBefore:  now.UnixMilli(),
		IssuedAt:   now.UnixMilli(),
		ID:         uuid.New().String(),
		LoggedinIn: user.LastLoggedinIn,
	}
}

func NewJsonTokenClaims(token *Token) *JsonTokenClaims {
	return &JsonTokenClaims{
		RegisteredClaims: &jwt.RegisteredClaims{
			Issuer:    token.Issuer,
			Subject:   strconv.Itoa(token.Subject),
			ExpiresAt: jwt.NewNumericDate(time.UnixMilli(token.ExpiresAt)),
			NotBefore: jwt.NewNumericDate(time.UnixMilli(token.NotBefore)),
			IssuedAt:  jwt.NewNumericDate(time.UnixMilli(token.IssuedAt)),
			ID:        token.ID,
		},
		UserName:   token.UserName,
		UserMail:   token.UserMail,
		LoggedinIn: token.LoggedinIn,
	}
}

func NewReqValidateToken() *ReqValidateToken {
	return &ReqValidateToken{
		ReqValidateTokenMeta: &ReqValidateTokenMeta{},
	}
}

//	func NewReqUpdateUserInfo(tk *TokenMeta) *ReqUpdateUserInfo {
//		return &ReqUpdateUserInfo{
//			ReqUpdateUserInfoMeta: NewReqUpdateUserInfoMeta(tk),
//		}
//	}
func NewReqUpdateUserInfo(tk *TokenMeta) *ReqUpdateUserInfo {
	return &ReqUpdateUserInfo{
		ReqUpdateUserInfoMeta: new(ReqUpdateUserInfoMeta),
	}
}

func (ruuim *ReqUpdateUserInfoMeta) SetReqUserInfoMeta(tk *TokenMeta) {
	ruuim.UserId = tk.Subject
	ruuim.UpdatedAt = time.Now().UnixMilli()
	ruuim.UpdatedIn = tk.LoggedinIn
	ruuim.JsonTokenId = tk.ID
}

// func NewReqUpdateUserInfoMeta(tk *TokenMeta) *ReqUpdateUserInfoMeta {
// 	return &ReqUpdateUserInfoMeta{
// 		UserId:      tk.Subject,
// 		UpdatedAt:   time.Now().UnixMilli(),
// 		UpdatedIn:   tk.LoggedinIn,
// 		JsonTokenId: tk.ID,
// 	}
// }

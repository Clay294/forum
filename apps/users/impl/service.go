package impl

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Clay294/forum/apps/users"

	"github.com/Clay294/forum/common"
	"github.com/Clay294/forum/rsakeys"

	"github.com/Clay294/forum/exception"
	"github.com/Clay294/forum/flog"
)

func (sc *serviceController) CreateUser(ctx context.Context, rcu *users.ReqCreateUser) error {
	// 1.验证请求参数的合法性
	err := rcu.Validate()

	if err != nil {
		return exception.NewErrInvalidRequest(err)
	}

	// 2.对密码进行哈希
	hp, err := HashPassword(rcu.Password)
	if err != nil {
		flog.Flogger().Error().Str("request_url", users.CreateUserUrl).Str("request_url", users.CreateUserMethod).Msgf("hashing password failed: %s", err)
		return exception.NewErrInternalServerError(fmt.Errorf("创建失败"))
	}
	rcu.Password = hp

	// 3.将user存入数据库
	err = sc.MySQLSaveUser(ctx, rcu)
	if err != nil {
		flog.Flogger().Error().Str("request_url", users.CreateUserUrl).Str("request_mehtod", users.CreateUserMethod).Msgf("creating user failed: %s", err)

		var feErr *exception.ForumException
		errors.As(err, &feErr)

		if feErr.ForumExceptionCode == exception.ErrNameExists {
			return exception.NewErrNameExists(fmt.Errorf("创建失败"))
		}

		if feErr.ForumExceptionCode == exception.ErrMailExists {
			return exception.NewErrMailExists(fmt.Errorf("创建失败"))
		}

		return exception.NewErrMySQLInternalError(fmt.Errorf("创建失败"))
	}
	return nil
}

func (sc *serviceController) GetUserInfo(ctx context.Context, rgui *users.ReqGetUserInfo) (*users.User, error) {
	// 1.验证请求参数的合法性
	err := rgui.Validate()

	if err != nil {
		return nil, exception.NewErrGetUserInfo(err)
	}

	// 2.从数据库中获取user信息
	user, err := sc.MySQLGetUserInfo(ctx, rgui)

	if err != nil {
		flog.Flogger().Error().Str("unit", sc.Name()).Str("request_url", users.GetUserInfoUrl).Str("request_mehtod", users.GetUserInfoMethod).Msgf("getting user info failed: %s", err)

		var feErr *exception.ForumException
		errors.As(err, &feErr)
		if feErr.ForumExceptionCode == exception.ErrUserNotFound {
			return nil, exception.NewErrUserNotFound(fmt.Errorf("用户不存在"))
		}

		return nil, exception.NewErrMySQLInternalError(fmt.Errorf("获取失败"))
	}

	return user, nil
}

func (sc *serviceController) Login(ctx context.Context, rl *users.ReqLogin) (*users.Token, error) {
	// 1.验证请求参数的合法性
	err := rl.Validate()

	if err != nil {
		return nil, exception.NewErrInvalidRequest(err)
	}

	// 2.查询user
	user, err := sc.MySQLQueryUser(ctx, rl)

	if err != nil {
		flog.Flogger().Error().Str("request_url", users.LoginUrl).Str("request_mehtod", users.LoginMethod).Msgf("getting user failed: %s", err)
		var feErr *exception.ForumException
		errors.As(err, &feErr)
		if feErr.ForumExceptionCode == exception.ErrUserNotFound {
			return nil, exception.NewErrUserNotFound(fmt.Errorf("用户不存在"))
		}

		return nil, exception.NewErrMySQLInternalError(fmt.Errorf("登陆失败"))
	}

	// 3.验证密码
	err = ValidatePassword(user.Password, rl.Password)

	if err != nil {
		flog.Flogger().Error().Str("unit", sc.Name()).Str("request_url", users.LoginUrl).Str("request_method", users.LoginMethod).Msgf("wrong password: %s", err)
		return nil, exception.NewErrUnauthorized(fmt.Errorf("密码错误"))
	}

	// 4.更新查询到的用户的最近一次登录信息
	user.LastLoggedinAt, user.LastLoggedinIn = rl.LoggedinAt, rl.LoggedinIn

	// 5.随机获取一个rsa密钥字符串的uuid
	uuid := rsakeys.GlobalRSAKeyPairsUUIDs()[common.Rander.Intn(len(rsakeys.GlobalRSAKeyPairsUUIDs()))]

	// 6.获取rsa私钥
	pri, err := sc.GetRSAPrivateKey(ctx, uuid)
	if err != nil {
		flog.Flogger().Error().Str("unit", sc.Name()).Str("request_url", users.LoginUrl).Str("request_method", users.LoginMethod).Msgf(err.Error())
		return nil, exception.NewErrInternalServerError(fmt.Errorf("登陆失败"))
	}

	// 5.生成token（随机rsa私钥签名）
	token, err := sc.SetJsonToken(user, pri, uuid)

	if err != nil {
		flog.Flogger().Error().Str("unit", sc.Name()).Str("request_url", users.LoginUrl).Str("request_method", users.LoginMethod).Msgf("generating json token failed: %s", err)
		return nil, exception.NewErrInternalServerError(fmt.Errorf("登陆失败"))
	}

	// 6.更新数据库中用户的最近一次登录信息并保存token到数据库
	err = sc.MySQLLogin(ctx, token)
	if err != nil {
		flog.Flogger().Error().Str("request_url", users.LoginUrl).Str("request_method", users.LoginMethod).Msgf("updating user lastloggedin info and saving token failed: %s", err)
		return nil, exception.NewErrMySQLInternalError(fmt.Errorf("登陆失败"))
	}

	// 7.返回token
	return token, nil
}

func (sc *serviceController) Logout(ctx context.Context, rl *users.ReqLogout) error {
	err := MySQLSaveTokenBlacklist(sc.gdbForum, ctx, rl.JsonTokenId)
	if err != nil {
		flog.Flogger().Error().Str("unit", sc.Name()).Str("request_url", users.LogoutUrl).Str("request_method", users.LogoutMethod).Msgf("logout failed: %s", err)

		return exception.NewErrMySQLInternalError(fmt.Errorf("登出失败"))
	}

	return nil
}

func (sc *serviceController) ValidateToken(ctx context.Context, rvt *users.ReqValidateToken) (*users.TokenMeta, error) {
	err := rvt.Validate()

	if err != nil {
		return nil, exception.NewErrInvalidRequest(err)
	}

	jt, err := sc.ParseJsonToken(ctx, rvt)
	if err != nil {
		flog.Flogger().Error().Str("request_url", users.ValidateTokenUrl).Str("requst_method", users.ValidateTokenMethod).Msgf("parsing token failed: %s", err)

		var feErr *exception.ForumException
		errors.As(err, &feErr)
		if feErr.ForumExceptionCode == http.StatusUnauthorized {
			return nil, exception.NewErrUnauthorized(fmt.Errorf("请重新登陆"))
		}

		return nil, exception.NewErrInternalServerError(fmt.Errorf("请重新登录"))
	}

	tokenMeta, err := sc.ValidateJsonTokenClaims(ctx, jt, rvt)
	if err != nil {
		flog.Flogger().Error().Str("request_url", users.ValidateTokenUrl).Str("requst_method", users.ValidateTokenMethod).Msgf("validating the json token payload is invalid: %s", err)
		var feErr *exception.ForumException
		errors.As(err, &feErr)

		if feErr.ForumExceptionCode == http.StatusUnauthorized {
			return nil, exception.NewErrUnauthorized(fmt.Errorf("请重新登录"))
		}

		return nil, exception.NewErrMySQLInternalError(fmt.Errorf("请重新登录"))
	}

	return tokenMeta, nil
}

func (sc *serviceController) UpdateUser(ctx context.Context, ruui *users.ReqUpdateUserInfo) error {
	err := ruui.Validate()
	if err != nil {
		return exception.NewErrInvalidRequest(err)
	}

	// 2.如果是修改密码，对修改后的密码进行哈希
	if ruui.Password != "" {
		hp, err := HashPassword(ruui.Password)
		if err != nil {
			flog.Flogger().Error().Str("request_url", users.CreateUserUrl).Str("request_url", users.CreateUserMethod).Msgf("hashing password failed: %s", err)
			return exception.NewErrInternalServerError(fmt.Errorf("修改失败"))
		}

		ruui.Password = hp
	}

	// 3.更新数据库中的user数据
	err = sc.MySQLUpdateUser(ctx, ruui)
	if err != nil {
		flog.Flogger().Error().Str("request_url", users.UpdateUserInfoUrl).Str("request_mehtod", users.UpdateUserInfoMethod).Msgf("updating user failed: %s", err)

		var feErr *exception.ForumException
		errors.As(err, &feErr)
		if feErr.ForumExceptionCode == exception.ErrNameExists {
			return exception.NewErrNameExists(fmt.Errorf("修改失败"))
		}

		if feErr.ForumExceptionCode == exception.ErrMailExists {
			return exception.NewErrMailExists(fmt.Errorf("修改失败"))
		}

		return exception.NewErrMySQLInternalError(fmt.Errorf("修改失败"))
	}

	return nil
}

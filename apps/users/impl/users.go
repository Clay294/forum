package impl

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Clay294/forum/apps/users"

	"github.com/Clay294/forum/exception"
	"github.com/Clay294/forum/ioc"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Clay294/forum/rsakeys"
	"golang.org/x/crypto/bcrypt"
)

//func HashPassword(req any) error {
//	hashErr := fmt.Errorf("哈希密码失败")
//
//	switch exact := req.(type) {
//	case *users.ReqCreateUser:
//		hp, err := bcrypt.GenerateFromPassword([]byte(exact.Password), 10)
//		if err != nil {
//			flog.Flogger().Error().Str("unit", users.UnitName).Str("request_url", users.CreateUserUrl).Msgf("hashing password failed: %s", err)
//			return exception.NewErrCreateUser(fmt.Errorf("创建失败: %s", hashErr))
//		}
//		exact.Password = string(hp)
//		return nil
//	case *users.ReqUpdateUserInfo:
//		hp, err := bcrypt.GenerateFromPassword([]byte(exact.Password), 10)
//		if err != nil {
//			flog.Flogger().Error().Str("unit", users.UnitName).Str("request_url", users.UpdateUserInfoUrl).Str("request_method", users.UpdateUserInfoMethod).Msgf("hashing password failed: %s", err)
//			//return exception.NewErrUpdateUserInfo(fmt.Errorf("修改失败"))
//			return exception.NewErrUpdateUserInfo(fmt.Errorf("修改失败: %s", hashErr))
//		}
//		exact.Password = string(hp)
//		return nil
//	default:
//		flog.Flogger().Error().Str("unit", users.UnitName).Msgf("hashing password failed: %s", "invalid request")
//		return exception.NewErrInvalidRequest(fmt.Errorf("操作失败: %s", "请求无效"))
//	}
//}

func HashPassword(p string) (string, error) {
	hp, err := bcrypt.GenerateFromPassword([]byte(p), 10)
	if err != nil {
		return "", err
	}
	return string(hp), nil
}

func ValidatePassword(hp, p string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hp), []byte(p))
	if err != nil {
		return err
	}

	return nil
}

func (sc *serviceController) GetRSAPrivateKey(ctx context.Context, uuid string) (*rsa.PrivateKey, error) {
	// 1.从ioc中获取keys单元的业务逻辑控制器
	res := ioc.Controllers().GetServiceController(rsakeys.UnitName)

	// 2.将获取到的结果进行断言为keys.Service接口
	service := res.(rsakeys.Service)

	// 3.调用keys.Service接口中的GetRSAPrivateKey方法
	pri, err := service.GetRSAPrivateKey(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return pri, nil
}

func (sc *serviceController) SetJsonToken(user *users.User, pri *rsa.PrivateKey, uuid string) (*users.Token, error) {
	// 1.根据user创建新token
	token := users.NewToken(user)

	// 2.根据token创建用于生成jwt的payload
	jtClaims := users.NewJsonTokenClaims(token)

	// 3.创建带payload的jwt，指定签名方式rsa256
	jt := jwt.NewWithClaims(jwt.SigningMethodRS256, jtClaims)

	// 4.将rsa密钥对儿的uuid放入jwt的header
	jt.Header["kid"] = uuid

	// 5.使用rsa私钥签名
	jtRaw, err := jt.SignedString(pri)
	if err != nil {
		return nil, err
	}

	// 6.使用生成的jwt更新token的JsonToken字段
	token.JsonToken = jtRaw

	// 7.返回生成的jwt
	return token, nil
}

func (sc *serviceController) GetRSAPublicKey(ctx context.Context, uuid string) (*rsa.PublicKey, error) {
	res := ioc.Controllers().GetServiceController(rsakeys.UnitName)

	service := res.(rsakeys.Service)

	pub, err := service.GetRSAPublicKey(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return pub, nil
}

func (sc *serviceController) ParseJsonToken(ctx context.Context, rvt *users.ReqValidateToken) (*jwt.Token, error) {
	// 1.解析json token的签名并验证基本信息
	jt, err := jwt.ParseWithClaims(
		rvt.JsonToken,
		// 提供&users.JsonTokenClaims{}结构体实例，用于解析时验证json token是否符合其相应的格式
		&users.JsonTokenClaims{},
		// 获取解析签名所需的公钥函数
		func(token *jwt.Token) (interface{}, error) {
			// 1. 从json token的header中获取kid
			res, ok := token.Header["kid"]
			if !ok {
				return nil, exception.NewErrUnauthorized(fmt.Errorf("the kid field is not in json token header"))
			}

			kid, ok := res.(string)
			if !ok || res == "" {
				return nil, exception.NewErrUnauthorized(fmt.Errorf("the kid in json token header is invalid"))
			}

			// kid = "130182c6-40ea-4e61-954f-71fd7007" // TODO debug
			//2. 获取kid对应的公钥
			pub, err := sc.GetRSAPublicKey(ctx, kid)
			if err != nil {
				return nil, exception.NewErrInternalServerError(err)
			}

			// 3.返回解析签名所需的公钥
			return pub, nil
		},
	)

	if err != nil {
		if feErr := (*exception.ForumException)(nil); errors.As(err, &feErr) {
			return nil, feErr
		}

		return nil, exception.NewErrUnauthorized(err)
	}

	// 2.验证token
	if !jt.Valid {
		return nil, exception.NewErrUnauthorized(err)
	}

	return jt, nil
}

func (sc *serviceController) ValidateJsonTokenClaims(ctx context.Context, jt *jwt.Token, rvt *users.ReqValidateToken) (*users.TokenMeta, error) {
	// 1.将解析后的json token的payload部分断言为*users.JsonTokenClaims
	jtClaims, ok := jt.Claims.(*users.JsonTokenClaims)

	if !ok {
		return nil, exception.NewErrUnauthorized(fmt.Errorf("the payload of json token is not in defined format"))
	}

	//  2.根据json token的payload中的jti查找其是否在tokens_blcaklist中
	err := sc.MySQLIsBlacklistToken(ctx, jtClaims.ID)
	if err != nil {
		var feErr *exception.ForumException
		errors.As(err, &feErr)

		if feErr.ForumExceptionCode == http.StatusUnauthorized {
			return nil, exception.NewErrUnauthorized(err)
		}

		return nil, exception.NewErrMySQLInternalError(err)
	}

	// 3.根据json token的payload中的jti从数据库中查找token meta
	tokenMeta, err := sc.MySQLGetTokenMeta(ctx, jtClaims.ID)

	if err != nil {
		var feErr *exception.ForumException
		errors.As(err, &feErr)

		if feErr.ForumExceptionCode == exception.ErrMySQLInternalError {
			return nil, exception.NewErrMySQLInternalError(err)
		}

		return nil, exception.NewErrUnauthorized(err)
	}

	// 4.验证json token中的payload是否合法
	if jtClaims.Issuer != tokenMeta.Issuer || jtClaims.Subject != strconv.Itoa(tokenMeta.Subject) || jtClaims.UserName != tokenMeta.UserName || jtClaims.UserMail != tokenMeta.UserMail || jtClaims.LoggedinIn != tokenMeta.LoggedinIn || jtClaims.LoggedinIn != rvt.ClientIP {
		return nil, exception.NewErrUnauthorized(fmt.Errorf("the json token payload is invalid"))
	}

	return tokenMeta, nil
}

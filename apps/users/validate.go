package users

import (
	"errors"
	"fmt"

	"github.com/Clay294/forum/flog"
	"github.com/go-playground/validator/v10"
)

func (rcu *ReqCreateUser) Validate() error {
	err := ValidateReqToUsers(rcu)
	if err != nil {
		var ivErr *validator.InvalidValidationError
		if errors.As(err, &ivErr) {
			flog.Flogger().Error().Str("unit", UnitName).Str("request_url", CreateUserUrl).Str("request_url", CreateUserMethod).Msgf("an internal error occurred during validating the request: %s", ivErr)
			return fmt.Errorf("校验创建用户请求参数时发生内部错误")
		}

		var vERRs validator.ValidationErrors
		if errors.As(err, &vERRs) {
			for _, err := range vERRs {
				switch err.StructField() {
				case "Name":
					switch err.Tag() {
					case "required":
						flog.Flogger().Error().Str("unit", UnitName).Str("request_url", CreateUserUrl).Str("request_method", CreateUserMethod).Msgf("the name filed is empty: %s", err)
						return fmt.Errorf("请输入用户名")
					case "validUserInfo":
						flog.Flogger().Error().Str("unit", UnitName).Str("request_url", CreateUserUrl).Str("request_method", CreateUserMethod).Msgf("the name is invalid: %s", err)
						return fmt.Errorf("大小写英文字母，数字，汉字，下划线")
					case "min", "max":
						flog.Flogger().Error().Str("unit", UnitName).Str("request_url", CreateUserUrl).Str("request_method", CreateUserMethod).Msgf("the name length is invalid: %s", err)
						return fmt.Errorf("1-25个字符")
					}
				case "Password":
					switch err.Tag() {
					case "required":
						flog.Flogger().Error().Str("unit", UnitName).Str("request_url", CreateUserUrl).Str("request_method", CreateUserMethod).Msgf("the password filed is empty: %s", err)
						return fmt.Errorf("请输入密码")
					case "validUserInfo":
						flog.Flogger().Error().Str("unit", UnitName).Str("request_url", CreateUserUrl).Str("request_method", CreateUserMethod).Msgf("password chracters from user creation request are invalid: %s", err)
						return fmt.Errorf("大小写英文字母，数字，特殊字符，必须包含大写或小写英文字母")
					case "min", "max":
						flog.Flogger().Error().Str("unit", UnitName).Str("request_url", CreateUserUrl).Str("request_method", CreateUserMethod).Msgf("the password length is invalid: %s", err)
						return fmt.Errorf("1-25个字符")
					}
				case "Mail":
					switch err.Tag() {
					case "required":
						flog.Flogger().Error().Str("unit", UnitName).Str("request_url", CreateUserUrl).Str("request_method", CreateUserMethod).Msgf("the mail filed is empty: %s", err)
						return fmt.Errorf("请输入邮箱")
					case "validUserInfo":
						flog.Flogger().Error().Str("unit", UnitName).Str("request_url", CreateUserUrl).Str("request_method", CreateUserMethod).Msgf("the mail is invalid: %s", err)
						return fmt.Errorf("qq，网易，搜狐，新浪邮箱")
					}
				}
			}
		}
	}
	return nil
}

// TODO
func (rgui *ReqGetUserInfo) Validate() error {
	err := ValidateReqToUsers(rgui)

	var ivErr *validator.InvalidValidationError
	if errors.As(err, &ivErr) {
		flog.Flogger().Error().Str("unit", UnitName).Str("request_url", GetUserInfoUrl).Str("request_method", GetUserInfoMethod).Msgf("an internal error occurred during validating the request: %s", ivErr)
		return fmt.Errorf("校验获取用户信息请求参数时发生内部错误")
	}

	var vERRs validator.ValidationErrors
	if errors.As(err, &vERRs) {
		for _, err := range vERRs {
			switch err.Tag() {
			case "required":
				flog.Flogger().Error().Str("unit", UnitName).Str("request_url", GetUserInfoUrl).Str("request_method", GetUserInfoMethod).Msgf("the id filed is empty: %s", err)
				return fmt.Errorf("请输入用户id")
			case "number":
				flog.Flogger().Error().Str("unit", UnitName).Str("request_url", GetUserInfoUrl).Str("request_method", GetUserInfoMethod).Msgf("the id is invalid: %s", err)
				return fmt.Errorf("用户id无效")
			}
		}
	}
	return nil
}

func (rl *ReqLogin) Validate() error {
	err := ValidateReqToUsers(rl)

	if err != nil {

		var ivErr *validator.InvalidValidationError

		if errors.As(err, &ivErr) {
			flog.Flogger().Error().Str("unit", UnitName).Str("request_url", LoginUrl).Str("request_method", LoginMethod).Msgf("an internal error occurred during validating the request: %s", ivErr)
			return fmt.Errorf("校验登录请求参数时发生内部错误")
		}

		var vERRs validator.ValidationErrors

		if errors.As(err, &vERRs) {
			for _, err := range vERRs {
				switch err.StructField() {
				case "Credential":
					flog.Flogger().Error().Str("unit", UnitName).Str("request_url", LoginUrl).Str("request_method", LoginMethod).Msgf("the name/mail filed is empty: %s", err)
					return fmt.Errorf("请输入登录用户名/邮箱")
				case "Password":
					flog.Flogger().Error().Str("unit", UnitName).Str("request_url", LoginUrl).Str("request_method", LoginMethod).Msgf("the password is empty: %s", err)
					return fmt.Errorf("请输入密码")
				}
			}

		}
	}
	return nil
}

func (rvt *ReqValidateToken) Validate() error {
	err := ValidateReqToUsers(rvt)
	if err != nil {
		var ivErr *validator.InvalidValidationError
		if errors.As(err, &ivErr) {
			flog.Flogger().Error().Str("unit", UnitName).Str("request_url", ValidateTokenUrl).Str("request_method", ValidateTokenMethod).Msgf("an internal error occurred during validating the request: %s", ivErr)
			return fmt.Errorf("校验验证token请求参数时发生内部错误")
		}
		var vERRs validator.ValidationErrors
		if errors.As(err, &vERRs) {
			for _, err := range vERRs {
				switch err.Tag() {
				case "required":
					flog.Flogger().Error().Str("unit", UnitName).Str("request_url", ValidateTokenUrl).Str("request_method", ValidateTokenMethod).Msgf("the json token field is empty: %s", ivErr)
					return fmt.Errorf("请先登录")
				case "jwt":
					flog.Flogger().Error().Str("unit", UnitName).Str("request_url", ValidateTokenUrl).Str("request_method", ValidateTokenMethod).Msgf("the json token is invalid: %s", ivErr)
					return fmt.Errorf("请先登录")
				}
			}
		}
	}
	return nil
}

func (ruui *ReqUpdateUserInfo) Validate() error {
	err := ValidateReqToUsers(ruui)

	if err != nil {
		var ivErr *validator.InvalidValidationError
		if errors.As(err, &ivErr) {
			flog.Flogger().Error().Str("unit", UnitName).Str("request_url", UpdateUserInfoUrl).Str("request_method", UpdateUserInfoMethod).Msgf("an internal error occurred during validating the request: %s", ivErr)
			return fmt.Errorf("校验更新用户信息请求时发生内部错误")
		}

		var vERRs validator.ValidationErrors
		if errors.As(err, &vERRs) {
			for _, err := range vERRs {
				switch err.StructField() {
				case "Name":
					switch err.Tag() {
					case "required_without_all":
						flog.Flogger().Error().Str("unit", UnitName).Str("request_url", UpdateUserInfoUrl).Str("request_method", UpdateUserInfoMethod).Msgf("the name filed is empty: %s", ivErr)
						return fmt.Errorf("请输入新用户名")
					case "excluded_with":
						flog.Flogger().Error().Str("unit", UnitName).Str("request_url", UpdateUserInfoUrl).Str("request_method", UpdateUserInfoMethod).Msgf("addtional user item has been selected to update: %s", ivErr)
						return fmt.Errorf("已选择其它修改项")
					case "min", "max":
						flog.Flogger().Error().Str("unit", UnitName).Str("request_url", UpdateUserInfoUrl).Str("request_method", UpdateUserInfoMethod).Msgf("the name length is invalid: %s", err)
						return fmt.Errorf("1-25个字符")
					case "validUserInfo":
						flog.Flogger().Error().Str("unit", UnitName).Str("request_url", UpdateUserInfoUrl).Str("request_method", UpdateUserInfoMethod).Msgf(" the name chracters are invalid: %s", err)
						return fmt.Errorf("大小写英文字母，数字，汉字，下划线")
					}
				case "Mail":
					switch err.Tag() {
					case "required_without_all":
						flog.Flogger().Error().Str("unit", UnitName).Str("request_url", UpdateUserInfoUrl).Str("request_method", UpdateUserInfoMethod).Msgf("the mail filed is empty: %s", ivErr)
						return fmt.Errorf("请输入新邮箱")
					case "excluded_with":
						flog.Flogger().Error().Str("unit", UnitName).Str("request_url", UpdateUserInfoUrl).Str("request_method", UpdateUserInfoMethod).Msgf("addtional user item has been selected to update: %s", ivErr)
						return fmt.Errorf("已选择其它修改项")
					case "validUserInfo":
						flog.Flogger().Error().Str("unit", UnitName).Str("request_url", UpdateUserInfoUrl).Str("request_method", UpdateUserInfoMethod).Msgf(" the mail chracters are invalid: %s", err)
						return fmt.Errorf("qq，网易，搜狐，新浪邮箱")

					}
				case "Password":
					switch err.Tag() {
					case "required_without_all":
						flog.Flogger().Error().Str("unit", UnitName).Str("request_url", UpdateUserInfoUrl).Str("request_method", UpdateUserInfoMethod).Msgf("the password field is empty: %s", ivErr)
						return fmt.Errorf("请输入新密码")
					case "excluded_with":
						flog.Flogger().Error().Str("unit", UnitName).Str("request_url", UpdateUserInfoUrl).Str("request_method", UpdateUserInfoMethod).Msgf("addtional user item has been selected to update: %s", ivErr)
						return fmt.Errorf("已选择其它修改项")
					case "min", "max":
						flog.Flogger().Error().Str("unit", UnitName).Str("request_url", UpdateUserInfoUrl).Str("request_method", UpdateUserInfoMethod).Msgf("the password length is invalid: %s", err)
						return fmt.Errorf("1-25个字符")
					case "validUserInfo":
						flog.Flogger().Error().Str("unit", UnitName).Str("request_url", UpdateUserInfoUrl).Str("request_method", UpdateUserInfoMethod).Msgf(" the password chracters are invalid: %s", err)
						return fmt.Errorf("大小写英文字母，数字，特殊字符，必须包含大写或小写英文字母")
					}
				}
			}
		}
	}
	return nil
}

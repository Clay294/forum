package users

import (
	"github.com/Clay294/forum/flog"
	"github.com/dlclark/regexp2"
	"github.com/go-playground/validator/v10"
)

const (
	namePattern     = `^(?=.*[a-zA-Z])[a-zA-Z0-9\u4e00-\u9fa5_]+$`
	mailPattern     = `^.+@(qq|sohu|163|162|sina).com$`
	passwordPattern = `^(?=.*[a-zA-Z])[~!@#$%^&*()\-_+=<>?:"{}|,.;'\[\]a-zA-Z0-9]+$`
)

var nRE, mRE, pRE *regexp2.Regexp

var usersValidator = validator.New()

func ValidateReqToUsers(req any) error {
	err := usersValidator.Struct(req)
	return err
}

func validUserInfoFunc(fl validator.FieldLevel) bool {
	switch fl.FieldName() {
	//case "Name", "UserName":
	case "Name":
		isValid, err := nRE.MatchString(fl.Field().String())
		if err != nil {
			flog.Flogger().Error().Str("unit", UnitName).Msgf("matching the mail by regular expression timeout: %s", err)
			return false
		}

		if !isValid {
			return false
		}

		return true
	//case "Mail", "UserMail":
	case "Mail":
		isValid, err := mRE.MatchString(fl.Field().String())

		if err != nil {
			flog.Flogger().Error().Msgf("matching the mail by regexp expression timeout: %s", err)
			return false
		}

		if !isValid {
			return false
		}

		return true
	//case "Password", "UserPassword":
	case "Password":
		isValid, err := pRE.MatchString(fl.Field().String())
		if err != nil {
			flog.Flogger().Error().Msgf("matching the password by regexp expression timeout: %s", err)
			return false
		}

		if !isValid {
			return false
		}

		return true
	default:
		return false
	}
}

func init() {
	var err error
	nRE, err = regexp2.Compile(namePattern, 0)
	if err != nil {
		flog.Flogger().Panic().Msgf("compilation of the regular expression for the username failed:%s", err)
	}

	mRE, err = regexp2.Compile(mailPattern, 0)
	if err != nil {
		flog.Flogger().Panic().Msgf("compilation of the regular expression for the mail failed:%s", err)
	}
	pRE, err = regexp2.Compile(passwordPattern, 0)

	if err != nil {
		flog.Flogger().Panic().Msgf("compilation of the regular expression for the password failed:%s", err)
	}

	err = usersValidator.RegisterValidation("validUserInfo", validUserInfoFunc)
	if err != nil {
		flog.Flogger().Panic().Msgf("initialization of the request parameter validator failed:%s", err)
	}
}

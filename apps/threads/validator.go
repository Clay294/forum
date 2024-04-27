package threads

import (
	"github.com/Clay294/forum/flog"
	"github.com/dlclark/regexp2"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

const (
	linkPattern = `(https?://)?(pan|yun)\.baidu\.com/s/[a-zA-Z0-9_-]+`
)

var threadsValidator = validator.New()

var lRE *regexp2.Regexp

func validateReqToThreads(req any) error {
	return threadsValidator.Struct(req)
}

func validLinkFunc(fl validator.FieldLevel) bool {
	isValid, err := lRE.MatchString(fl.Field().String())
	if err != nil {
		flog.Flogger().Error().Msgf("matching the link by regexp expression timeout: %s", err)
		return false
	}

	return isValid
}

func validMainSectionFunc(fl validator.FieldLevel) bool {
	_, ok := sectionTable[MAINSECTION(fl.Field().String())]
	return ok
}

func validSubSectionFunc(fl validator.FieldLevel) bool {
	parent := fl.Parent().Interface()

	switch exact := parent.(type) {
	case ReqCreateThreads:
		_, ok := sectionTable[exact.MainSection][exact.SubSection]
		return ok
	}

	return false
}

func init() {
	var err error

	lRE, err = regexp2.Compile(linkPattern, 0)
	if err != nil {
		log.Panic().Caller().Msgf("compiling link pattern failed: %s", err)
	}

	err = threadsValidator.RegisterValidation("validLink", validLinkFunc)
	if err != nil {
		log.Panic().Caller().Msgf("registering validLink validation failed: %s", err)
	}

	err = threadsValidator.RegisterValidation("validMainSection", validMainSectionFunc)
	if err != nil {
		log.Panic().Caller().Msgf("registering validMainSecton validation failed: %s", err)

	}

	err = threadsValidator.RegisterValidation("validSubSection", validSubSectionFunc)
	if err != nil {
		log.Panic().Caller().Msgf("registering validSubSecton validation failed: %s", err)
	}
}

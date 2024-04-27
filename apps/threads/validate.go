package threads

import (
	"errors"
	"fmt"

	"github.com/Clay294/forum/flog"
	"github.com/go-playground/validator/v10"
)

func (rct *ReqCreateThreads) Validate() error {
	err := validateReqToThreads(rct)
	if err != nil {
		var ivErr *validator.InvalidValidationError
		if errors.As(err, &ivErr) {
			flog.Flogger().Error().Str("request_url", CreateThreadUrl).Str("request_method", CreateThreadMethod).Msgf("an internal error occurred during validating the request")
			return fmt.Errorf("校验创建帖子请求参数时发生内部错误")
		}

		var vErrs validator.ValidationErrors
		if errors.As(err, &vErrs) {
			fmt.Println(vErrs, "%%%%%%") //TODO

			for _, err := range vErrs {
				switch err.StructField() {
				case "Title":
					flog.Flogger().Error().Str("unit", UnitName).Str("request_url", CreateThreadUrl).Str("request_method", CreateThreadMethod).Msgf("the title is empty")
					return fmt.Errorf("请输入标题")
				case "MainSection":
					switch err.Tag() {
					case "required":
						flog.Flogger().Error().Str("request_url", CreateThreadUrl).Str("request_method", CreateThreadMethod).Msgf("the mainsection is empty")
						return fmt.Errorf("请指定主版块")
					case "validMainSection":
						flog.Flogger().Error().Str("request_url", CreateThreadUrl).Str("request_method", CreateThreadMethod).Msgf("the mainsection is invalid")
						return fmt.Errorf("主版块无效")
					}
				case "SubSection":
					switch err.Tag() {
					case "required":
						flog.Flogger().Error().Str("request_url", CreateThreadUrl).Str("request_method", CreateThreadMethod).Msgf("the subsection is empty")
						return fmt.Errorf("请指定子版块")
					case "validMainSection":
						flog.Flogger().Error().Str("request_url", CreateThreadUrl).Str("request_method", CreateThreadMethod).Msgf("the subsection is invalid")
						return fmt.Errorf("子版块无效")
					}
				case "Text":
					flog.Flogger().Error().Str("request_url", CreateThreadUrl).Str("request_method", CreateThreadMethod).Msgf("the text is empty")
					return fmt.Errorf("请输入正文")
				case "Link":
					switch err.Tag() {
					case "required":
						flog.Flogger().Error().Str("request_url", CreateThreadUrl).Str("request_method", CreateThreadMethod).Msgf("the link is empty")
						return fmt.Errorf("请输入资源链接")
					case "validMainSection":
						flog.Flogger().Error().Str("request_url", CreateThreadUrl).Str("request_method", CreateThreadMethod).Msgf("the link is invalid")
						return fmt.Errorf("百度云链接")
					}
				case "LinkCode":
					flog.Flogger().Error().Str("request_url", CreateThreadUrl).Str("request_method", CreateThreadMethod).Msgf("the linkcode is empty")
					return fmt.Errorf("请输入资源链接提取码")
				case "UnzipPassword":
					flog.Flogger().Error().Str("request_url", CreateThreadUrl).Str("request_method", CreateThreadMethod).Msgf("the unzippassword is empty")
					return fmt.Errorf("请输入资源解压码")
				case "Price":
					switch err.Tag() {
					case "required":
						flog.Flogger().Error().Str("request_url", CreateThreadUrl).Str("request_method", CreateThreadMethod).Msgf("the price is empty")
						return fmt.Errorf("请输入定价")
					case "number":
						flog.Flogger().Error().Str("request_url", CreateThreadUrl).Str("request_method", CreateThreadMethod).Msgf("the price is invalid")
						return fmt.Errorf("定价无效")
					}
				// case "Tags":
				case "Status":
					flog.Flogger().Error().Str("request_url", CreateThreadUrl).Str("request_method", CreateThreadMethod).Msgf("the thread status is invalid")
					return fmt.Errorf("草稿/已发布")
				}
			}
		}
	}
	return nil
}

func (rsbmh *ReqSearchByMainHome) Validate() error {
	err := validateReqToThreads(rsbmh)
	if err != nil {
		var ivErr *validator.InvalidValidationError
		if errors.As(err, &ivErr) {
			flog.Flogger().Error().Str("unit", UnitName).Str("request_url", SearchByMainHomeUrl).Str("request_mehtod", SearchByMainHomeMethod).Msgf("an internal error occurred during validating the request")
			return fmt.Errorf("发生内部错误")
		}

		var vErrs validator.ValidationErrors
		if errors.As(err, &vErrs) {
			for _, err := range vErrs {
				switch err.StructField() {
				case "Keywords":
					flog.Flogger().Error().Str("unit", UnitName).Str("request_url", SearchByMainHomeUrl).Str("request_mehtod", SearchByMainHomeMethod).Msgf("searching by author has been selected")
					return fmt.Errorf("已指定按作者搜索")
				case "UserName":
					flog.Flogger().Error().Str("unit", UnitName).Str("request_url", SearchByMainHomeUrl).Str("request_mehtod", SearchByMainHomeMethod).Msgf("searching by keywords has been selected")
					return fmt.Errorf("已指定按关键字搜索")
				}
			}
		}
	}
	return nil
}

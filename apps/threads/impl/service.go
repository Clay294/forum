package impl

import (
	"context"
	"errors"
	"fmt"

	"github.com/Clay294/forum/apps/threads"
	"github.com/Clay294/forum/exception"
	"github.com/Clay294/forum/flog"
)

func (sc *serviceController) CreateThread(ctx context.Context, rct *threads.ReqCreateThreads) error {
	err := rct.Validate()
	if err != nil {
		return exception.NewErrInvalidRequest(err)
	}

	err = sc.MySQLSaveThread(ctx, rct)
	if err != nil {
		flog.Flogger().Error().Str("request_url", threads.CreateThreadUrl).Str("request_method", threads.CreateThreadMethod).Msgf("creating thread failed: %s", err)

		var feErr *exception.ForumException
		errors.As(err, &feErr)

		if feErr.ForumExceptionCode == exception.ErrTitleExists {
			return exception.NewErrTitleExists(fmt.Errorf("标题已存在"))
		}

		return exception.NewErrMySQLInternalError(fmt.Errorf("创建失败"))
	}

	return nil
}

func (sc *serviceController) SearchByMainHome(ctx context.Context, rsbmh *threads.ReqSearchByMainHome) (*threads.ThreadsList, error) {
	err := rsbmh.Validate()
	if err != nil {
		return nil, exception.NewErrInvalidRequest(err)
	}

	tl, err := sc.SearchByMainHome(ctx, rsbmh)
	if err != nil {
		flog.Flogger().Error().Str("unit", sc.Name()).Str("request_url", threads.SearchByMainHomeUrl).Str("request_method", threads.SearchByMainHomeMethod).Msgf("search by mainhome failed: %s", err)
		return nil, exception.NewErrMySQLInternalError(fmt.Errorf("搜索失败"))
	}

	return tl, nil
}

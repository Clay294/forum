package impl

import (
	"context"
	"errors"
	"fmt"

	"github.com/Clay294/forum/apps/threads"
	"github.com/Clay294/forum/exception"
	"github.com/Clay294/forum/flog"
	"gorm.io/gorm"
)

func (sc *serviceController) MySQLSaveThread(ctx context.Context, rct *threads.ReqCreateThreads) error {
	thread := threads.NewThread(rct)

	db := sc.gdbForum.WithContext(ctx)

	err := db.Model(thread).Select("Title").Where("title = ?", thread.Title).First(thread).Error
	if err == nil {
		flog.Flogger().Error().Str("database", "forum").Str("table", threads.UnitName).Msgf("the title already exists: %s", err)
		return exception.NewErrTitleExists(err)
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		flog.Flogger().Error().Str("database", "forum").Str("table", threads.UnitName).Msgf("querying threads failed: %s", err)
		return exception.NewErrMySQLInternalError(err)
	}

	err = db.Table(sc.Name()).Save(thread).Error
	if err != nil {
		flog.Flogger().Error().Str("database", "forum").Str("table", threads.UnitName).Msgf("saving thread failed: %s", err)
		return exception.NewErrMySQLInternalError(err)
	}

	return nil
}

func (sc *serviceController) MySQLSearchByMainHome(ctx context.Context, rsbmh *threads.ReqSearchByMainHome) (*threads.ThreadsList, error) {
	tl := threads.NewThreadsList()

	db := sc.gdbForum.WithContext(ctx).Table(threads.UnitName)

	if rsbmh.Keywords != "" {
		db.Where("title LIKE ?", fmt.Sprintf("%%%s%%", rsbmh.Keywords)).Or("text LIKE ?", fmt.Sprintf("%%%s%%", rsbmh.Keywords))
	}

	if rsbmh.UserName != "" {
		db.Where("user_name LIKE ?", fmt.Sprintf("%%%s%%", rsbmh.UserName))
	}

	err := db.Count(&tl.Total).Offset(rsbmh.PageSize * (rsbmh.PageNumber - 1)).Limit(rsbmh.PageSize).Find(&tl.List).Error
	if err != nil {
		flog.Flogger().Error().Str("database", "forum").Str("table", threads.UnitName).Msgf("querying threads failed: %s", err)
		return nil, err
	}

	return nil, nil
}

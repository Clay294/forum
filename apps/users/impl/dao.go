package impl

import (
	"context"
	"errors"
	"fmt"

	"github.com/Clay294/forum/apps/users"
	"github.com/Clay294/forum/common"
	"github.com/Clay294/forum/exception"
	"github.com/Clay294/forum/flog"
	"github.com/Clay294/forum/rsakeys"
	"gorm.io/gorm"
)

func (sc *serviceController) MySQLSaveUser(ctx context.Context, rcu *users.ReqCreateUser) error {
	user := users.NewUser(rcu)

	db := sc.gdbForum.WithContext(ctx)

	err := db.Model(user).Select("name").Where("name = ?", rcu.Name).First(user).Error

	if err == nil {
		flog.Flogger().Error().Str("unit", sc.Name()).Str("database", "forum").Str("table", sc.Name()).Msgf("the name already exists")
		return exception.NewErrNameExists(fmt.Errorf("the name already exists"))
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		flog.Flogger().Error().Str("unit", sc.Name()).Str("database", "forum").Str("table", sc.Name()).Msgf("querying name failed: %s", err)
		return exception.NewErrMySQLInternalError(fmt.Errorf("querying name failed: %s", err))
	}

	err = db.Model(user).Select("mail").Where("mail = ?", rcu.Mail).First(user).Error

	if err == nil {
		flog.Flogger().Error().Str("unit", sc.Name()).Str("database", "forum").Str("table", sc.Name()).Msgf("the mail already exists")
		return exception.NewErrMailExists(fmt.Errorf("the mail already exists"))
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		flog.Flogger().Error().Str("unit", sc.Name()).Str("request_url", users.CreateUserUrl).Str("request_method", users.CreateUserMethod).Str("database", "forum").Str("table", sc.Name()).Msgf("querying mail failed: %s", err)
		return exception.NewErrMySQLInternalError(fmt.Errorf("querying mail failed: %s", err))
	}

	err = db.Save(user).Error

	if err != nil {
		flog.Flogger().Error().Str("unit", sc.Name()).Str("database", "forum").Str("table", sc.Name()).Msgf("saving user failed: %s", err)
		return exception.NewErrMySQLInternalError(fmt.Errorf("saving user failed: %s", err))
	}

	return nil
}

func (sc *serviceController) MySQLGetUserInfo(ctx context.Context, rgui *users.ReqGetUserInfo) (*users.User, error) {
	user := new(users.User)

	db := sc.gdbForum.WithContext(ctx)

	err := db.
		Select("id, name, mail, coin, created_at, created_in, last_loggedin_at, last_loggedin_in").
		Where("id = ?", rgui.UserId).
		First(user).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			flog.Flogger().Error().Str("unit", sc.Name()).Str("database", "forum").Str("table", sc.Name()).Msgf("the user does not exist")
			return nil, exception.NewErrUserNotFound(fmt.Errorf("the user does not exist"))
		}

		flog.Flogger().Error().Str("unit", sc.Name()).Str("database", "forum").Str("table", sc.Name()).Msgf("querying user by id failed: %s", err)
		return nil, exception.NewErrMySQLInternalError(fmt.Errorf("querying user by id failed: %s", err))
	}
	return user, nil
}

func (sc *serviceController) MySQLQueryUser(ctx context.Context, rl *users.ReqLogin) (*users.User, error) {
	user := new(users.User)

	db := sc.gdbForum.WithContext(ctx).Table(users.UnitName)

	switch rl.LoginBy {
	case users.LoginByName:
		db.Where("name = ?", rl.Credential)
	case users.LoginByMail:
		db.Where("mail = ?", rl.Credential)
	}

	err := db.First(user).Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			flog.Flogger().Error().Str("unit", sc.Name()).Str("database", "forum").Str("table", sc.Name()).Msgf("the user does not exist")
			return nil, exception.NewErrUserNotFound(fmt.Errorf("the user does not exist"))
		}

		flog.Flogger().Error().Str("unit", sc.Name()).Str("database", "forum").Str("table", sc.Name()).Msgf("querying user failed: %s", err)
		return nil, exception.NewErrMySQLInternalError(fmt.Errorf("querying user failed: %s", err))
	}

	return user, nil
}

func (sc *serviceController) MySQLLogin(ctx context.Context, token *users.Token) error {
	err := sc.gdbForum.
		WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {
			err := MySQLUpdateUserLastLoggedIn(tx, ctx, token)
			if err != nil {
				return exception.NewErrMySQLInternalError(err)
			}

			err = MySQLSaveToken(tx, ctx, token)
			if err != nil {
				return exception.NewErrMySQLInternalError(err)
			}

			return nil
		})

	if err != nil {
		var feErr *exception.ForumException
		if errors.As(err, &feErr) {
			return feErr
		}

		flog.Flogger().Error().Str("unit", users.UnitName).Str("database", "forum").Strs("table", []string{sc.Name(), users.TokensTableName}).Bool("transaction", true).Msgf("login transaction excuting failed: %s", err)
		return err
	}
	return nil
}

func MySQLUpdateUserLastLoggedIn(gdb *gorm.DB, ctx context.Context, token *users.Token) error {
	err := gdb.
		WithContext(ctx).
		Model(&users.User{}).
		Where("id = ?", token.Subject).
		Omit("updated_at").
		Updates(map[string]interface{}{
			"last_loggedin_at": token.IssuedAt,
			"last_loggedin_in": token.LoggedinIn,
		},
		).
		Error
	if err != nil {
		flog.Flogger().Error().Str("unit", users.UnitName).Str("database", "forum").Str("table", users.UnitName).Msgf("updating user lastloggedat and lastlggedinin info failed: %s", err)
		return fmt.Errorf("updating user lastloggedat and lastlggedinin info failed: %s", err)
	}

	return nil
}

func MySQLSaveToken(gdb *gorm.DB, ctx context.Context, token *users.Token) error {
	err := gdb.WithContext(ctx).Model(token).Save(token).Error

	if err != nil {
		flog.Flogger().Error().Str("unit", users.UnitName).Str("database", "forum").Str("table", users.UnitName).Msgf("saving token failed: %s", err)
		return fmt.Errorf("saving token failed: %s", err)
	}

	return nil
}

func (sc *serviceController) MySQLIsBlacklistToken(ctx context.Context, uuid string) error {
	tkInBlacklist := new(users.TokenInBlacklist)
	tkInBlacklist.UUID = uuid

	err := sc.gdbForum.WithContext(ctx).Model(tkInBlacklist).Where(tkInBlacklist).First(tkInBlacklist).Error
	if err == nil {
		flog.Flogger().Error().Str("unit", sc.Name()).Str("database", "forum").Str("table", users.TokensBlacklistTableName).Msgf("the json token is in blacklist")
		return exception.NewErrUnauthorized(fmt.Errorf("the json token is in blacklist"))
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		flog.Flogger().Error().Str("unit", sc.Name()).Str("database", "forum").Str("table", users.TokensBlacklistTableName).Msgf("querying json token in tokens blacklist failed: %s", err)
		return exception.NewErrMySQLInternalError(fmt.Errorf("querying json token in tokens blacklist failed: %s", err))
	}

	return nil
}

func (sc *serviceController) MySQLGetTokenMeta(ctx context.Context, uuid string) (*users.TokenMeta, error) {

	tokenMeta := new(users.TokenMeta)

	err := sc.gdbForum.WithContext(ctx).Model(tokenMeta).Where("uuid = ?", uuid).First(tokenMeta).Error

	if err != nil {
		flog.Flogger().Error().Str("unit", sc.Name()).Str("database", "forum").Str("table", users.TokensTableName).Msgf("querying json token in tokens failed: %s", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewErrUnauthorized(fmt.Errorf("the json token is not exist"))
		}
		return nil, exception.NewErrMySQLInternalError(fmt.Errorf("querying json token in tokens failed: %s", err))
	}

	return tokenMeta, nil
}

func (sc *serviceController) MySQLUpdateUser(ctx context.Context, ruui *users.ReqUpdateUserInfo) error {
	err := sc.gdbForum.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			if ruui.Name != "" || ruui.Mail != "" {
				user, err := MySQLUpdateUserSpecifiedInfo(tx, ctx, ruui)
				if err != nil {
					var feErr *exception.ForumException
					errors.As(err, &feErr)
					if feErr.ForumExceptionCode == exception.ErrNameExists {
						return exception.NewErrNameExists(err)
					}

					if feErr.ForumExceptionCode == exception.ErrMailExists {
						return exception.NewErrMailExists(err)
					}

					return exception.NewErrMySQLInternalError(err)
				}

				uuid := rsakeys.GlobalRSAKeyPairsUUIDs()[common.Rander.Intn(len(rsakeys.GlobalRSAKeyPairsUUIDs()))]

				pri, err := sc.GetRSAPrivateKey(ctx, uuid)
				if err != nil {
					return exception.NewErrInternalServerError(err)
				}

				token, err := sc.SetJsonToken(user, pri, uuid)
				if err != nil {
					return exception.NewErrInternalServerError(err)
				}

				err = MySQLSaveToken(tx, ctx, token)
				if err != nil {
					return exception.NewErrMySQLInternalError(err)
				}
			}

			if ruui.Password != "" {
				_, err := MySQLUpdateUserSpecifiedInfo(tx, ctx, ruui)
				if err != nil {
					var feErr *exception.ForumException
					errors.As(err, &feErr)
					if feErr.ForumExceptionCode == exception.ErrNameExists {
						return exception.NewErrNameExists(err)
					}

					if feErr.ForumExceptionCode == exception.ErrMailExists {
						return exception.NewErrMailExists(err)
					}

					return exception.NewErrMySQLInternalError(err)
				}

				// err = MySQLSaveTokenBlacklist(tx, ctx, ruui)
				err = MySQLSaveTokenBlacklist(tx, ctx, ruui.JsonTokenId)
				if err != nil {
					return exception.NewErrMySQLInternalError(err)
				}
			}
			return nil
		},
	)
	if err != nil {
		var feErr *exception.ForumException

		if errors.As(err, &feErr) {
			return feErr
		}

		flog.Flogger().Error().Str("database", "forum").Strs("table", []string{sc.Name(), users.TokensTableName}).Bool("tracsaction", true).Msgf("updating specified user data and saving token to blacklistg failed: %s", err)
		return exception.NewErrMySQLInternalError(err)
	}

	return nil
}

func MySQLUpdateUserSpecifiedInfo(gdb *gorm.DB, ctx context.Context, ruui *users.ReqUpdateUserInfo) (*users.User, error) {
	user := new(users.User)
	db := gdb.WithContext(ctx)
	if ruui.Name != "" {
		err := db.Model(user).Where("name = ?", ruui.Name).First(user).Error
		if err == nil {
			flog.Flogger().Error().Str("unit", users.UnitName).Str("database", "forum").Str("table", users.UnitName).Msgf("the name already exists")
			return nil, exception.NewErrNameExists(fmt.Errorf("the name already exists"))
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			flog.Flogger().Error().Str("unit", users.UnitName).Str("database", "forum").Str("table", users.UnitName).Msgf("querying name failed: %s", err)
			return nil, exception.NewErrMySQLInternalError(fmt.Errorf("querying name failed: %s", err))
		}

		db = db.Model(user).Where("id = ?", ruui.UserId).Updates(map[string]interface{}{"name": ruui.Name, "updated_at": ruui.UpdatedAt})
	}

	if ruui.Mail != "" {
		err := db.Model(user).Where("mail = ?", ruui.Mail).First(user).Error
		if err == nil {
			flog.Flogger().Error().Str("unit", users.UnitName).Str("database", "forum").Str("table", users.UnitName).Msgf("the mail already exists")
			return nil, exception.NewErrNameExists(fmt.Errorf("the mail already exists"))
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			flog.Flogger().Error().Str("unit", users.UnitName).Str("database", "forum").Str("table", users.UnitName).Msgf("querying name failed: %s", err)
			return nil, exception.NewErrMySQLInternalError(fmt.Errorf("querying mail failed: %s", err))
		}

		db = db.Model(user).Where("id = ?", ruui.UserId).Updates(map[string]interface{}{"mail": ruui.Mail, "updated_at": ruui.UpdatedAt})

	}

	if ruui.Password != "" {
		db = db.Model(user).Where("id = ?", ruui.UserId).Updates(map[string]interface{}{"password": ruui.Password, "updated_at": ruui.UpdatedAt})
	}

	err := db.First(user).Error
	if err != nil {
		flog.Flogger().Error().Str("unit", users.UnitName).Str("database", "forum").Str("table", users.UnitName).Msgf("updating specified user data failed: %s", err)
		return nil, exception.NewErrMySQLInternalError(fmt.Errorf("updating user specified info failed: %s", err))
	}

	return user, nil
}

// func MySQLSaveTokenBlacklist(gdb *gorm.DB, ctx context.Context, ruui *users.ReqUpdateUserInfo) error {
// 	tkInBlcaklist := &users.TokenInBlacklist{UUID: ruui.JsonTokenId}

//		err := gdb.WithContext(ctx).Table(users.TokensBlacklistTableName).Save(tkInBlcaklist).Error
//		if err != nil {
//			flog.Flogger().Error().Str("database", "forum").Str("table", users.TokensBlacklistTableName).Msgf("saving token failed: %s", err)
//			return fmt.Errorf("saving token failed: %s", err)
//		}
//		return nil
//	}
func MySQLSaveTokenBlacklist(gdb *gorm.DB, ctx context.Context, uuid string) error {
	tkInBlcaklist := &users.TokenInBlacklist{UUID: uuid}

	err := gdb.WithContext(ctx).Table(users.TokensBlacklistTableName).Save(tkInBlcaklist).Error
	if err != nil {
		flog.Flogger().Error().Str("database", "forum").Str("table", users.TokensBlacklistTableName).Msgf("saving token failed: %s", err)
		return fmt.Errorf("saving token failed: %s", err)
	}
	return nil
}

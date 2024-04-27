package impl

import (
	"context"
	"fmt"
	"github.com/Clay294/forum/flog"

	"github.com/Clay294/forum/rsakeys"
)

func (sc *serviceController) MySQLGetPrivateKeyString(ctx context.Context, uuid string) (string, error) {
	rsaKB := new(rsakeys.RSAKeyPair)

	err := sc.gdbKeys.
		WithContext(ctx).
		Model(rsaKB).
		Select("private_key_string").
		Where("uuid = ?", uuid).
		First(rsaKB).
		Error
	if err != nil {
		flog.Flogger().Error().Str("unit", sc.Name()).Str("database", "keys").Str("table", sc.Name()).Msgf("querying rsa private key string faield: %s", err)
		return "", fmt.Errorf("querying rsa private key string faield: %s", err)
	}

	return rsaKB.PrivateKeyString, nil
}

func (sc *serviceController) MySQLGetRSAPublicKeyString(ctx context.Context, uuid string) (string, error) {
	rsaKB := new(rsakeys.RSAKeyPair)

	err := sc.gdbKeys.
		WithContext(ctx).
		Model(rsaKB).
		Select("public_key_string").
		Where("uuid = ?", uuid).
		First(rsaKB).
		Error

	if err != nil {
		flog.Flogger().Error().Str("unit", sc.Name()).Str("database", "keys").Str("table", sc.Name()).Msgf("querying rsa public key string faield: %s", err)
		return "", fmt.Errorf("querying rsa public key string faield: %s", err)
	}
	return rsaKB.PublicKeyString, nil
}

package impl

import (
	"context"
	"crypto/rsa"
	"fmt"
	"github.com/Clay294/forum/flog"
)

func (sc *serviceController) GetRSAPrivateKey(ctx context.Context, uuid string) (*rsa.PrivateKey, error) {
	priStr, err := sc.MySQLGetPrivateKeyString(ctx, uuid)

	if err != nil {
		flog.Flogger().Error().Str("unit", sc.Name()).Msgf("getting rsa private key string failed: %s", err)
		return nil, fmt.Errorf("get rsa private key string failed: %s", err)
	}

	pri, err := ParseRSAPrivateKeyString(priStr)

	if err != nil {
		flog.Flogger().Error().Str("unit", sc.Name()).Msgf("parsing rsa private key string failed: %s", err)
		return nil, fmt.Errorf("get rsa private key string failed: %s", err)
	}

	return pri, nil
}

func (sc *serviceController) GetRSAPublicKey(ctx context.Context, uuid string) (*rsa.PublicKey, error) {
	pubStr, err := sc.MySQLGetRSAPublicKeyString(ctx, uuid)

	if err != nil {
		flog.Flogger().Error().Str("unit", sc.Name()).Msgf("getting rsa public key string failed: %s", err)
		return nil, fmt.Errorf("get rsa public key string failed: %s", err)
	}

	pub, err := ParseRSAPublicKeyString(pubStr)

	if err != nil {
		flog.Flogger().Error().Str("unit", sc.Name()).Msgf("parsing rsa public key string failed: %s", err)
		return nil, fmt.Errorf("get rsa public key string failed: %s", err)
	}

	return pub, nil
}

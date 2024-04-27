package rsakeys

import (
	"context"
	"crypto/rsa"
)

type Service interface {
	GetRSAPrivateKey(context.Context, string) (*rsa.PrivateKey, error)
	GetRSAPublicKey(context.Context, string) (*rsa.PublicKey, error)
}

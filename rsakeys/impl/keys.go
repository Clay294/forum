package impl

import (
	"crypto/rsa"
	"encoding/base64"

	"github.com/golang-jwt/jwt/v5"
)

func ParseRSAPrivateKeyString(priStr string) (*rsa.PrivateKey, error) {
	priPem, err := base64.StdEncoding.DecodeString(priStr)
	if err != nil {
		return nil, err
	}

	pri, err := jwt.ParseRSAPrivateKeyFromPEM(priPem)

	if err != nil {
		return nil, err
	}

	return pri, nil
}

func ParseRSAPublicKeyString(pubStr string) (*rsa.PublicKey, error) {
	pubPem, err := base64.StdEncoding.DecodeString(pubStr)

	if err != nil {
		return nil, err
	}

	pub, err := jwt.ParseRSAPublicKeyFromPEM(pubPem)

	if err != nil {
		return nil, err
	}

	return pub, nil
}

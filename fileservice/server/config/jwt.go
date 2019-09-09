package config

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/dgrijalva/jwt-go"
)

type Jwt struct {
	Audience string `yaml:"audience"`
	Issuer   string `yaml:"issuer"`
	Subject  string `yaml:"subject"`
	Expire   int64  `yaml:"expire"`

	SignKey       string            `yaml:"sign_key"`
	signKeyByte   []byte            `yaml:"-"`
	SigningMethod string            `yaml:"signing_method"`
	signedMethod  jwt.SigningMethod `yaml:"-"`
}

func (p *Jwt) init() {
	p.signedMethod = jwt.GetSigningMethod(p.SigningMethod)
	if p.signedMethod == nil {
		panic(fmt.Sprintf("jwt not support SigningMethod: %s", p.SigningMethod))
	}

	if p.Expire <= 0 {
		p.Expire = 10
	}
	p.Expire = int64(time.Minute) * p.Expire

	p.signKeyByte = []byte(p.SignKey)
}

func (p *Jwt) newClaims(user string) jwt.StandardClaims {
	var now = time.Now().Unix()
	return jwt.StandardClaims{
		Audience:  p.Audience,
		ExpiresAt: now + p.Expire,
		Id:        user,
		IssuedAt:  now,
		Issuer:    p.Issuer,
		Subject:   p.Subject,
		NotBefore: now,
	}
}

func (p *Jwt) CreateToken(user string) (string, error) {
	var token = jwt.NewWithClaims(p.signedMethod, p.newClaims(user))

	return token.SignedString(p.signKeyByte)
}

func (p *Jwt) Verify(tokenString string) (string, error) {
	var token, err = jwt.ParseWithClaims(tokenString, jwt.StandardClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		claims, ok := token.Claims.(*jwt.StandardClaims)
		if !ok {
			return nil, errors.New("Unexpected Claims Type")
		}
		if claims.Subject != p.Subject {
			return nil, errors.New("unexpected Subject")
		}
		if err := claims.Valid(); err != nil {
			return nil, err
		}

		return p.signKeyByte, nil
	})
	if err != nil {
		return "", fmt.Errorf("token parse failed, err: %v", err)
	}

	if token == nil {
		return "", errors.New("token parse failed, token is nil")
	}

	return token.Claims.(*jwt.StandardClaims).Id, nil
}

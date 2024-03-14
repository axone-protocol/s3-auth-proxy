package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"time"
)

func (a *Authenticator) issueJwt() (string, error) {
	now := time.Now()
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Audience:  "", // TODO: put the authenticated service identifier here
		ExpiresAt: now.Add(5 * time.Minute).Unix(),
		Id:        uuid.New().String(),
		IssuedAt:  now.Unix(),
		Issuer:    a.serviceID,
		NotBefore: now.Unix(),
		Subject:   a.serviceID,
	}).SignedString(a.jwtSecretKey)
}

func (a *Authenticator) verifyJwt(raw string) error {
	token, err := jwt.Parse(raw, func(_ *jwt.Token) (interface{}, error) {
		return a.jwtSecretKey, nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}
	return nil
}

package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"time"
)

type Authenticator struct {
	jwtSecretKey    []byte
	node            string
	cognitariumAddr string
	serviceID       string
}

func New(jwtSecretKey []byte, node, cognitariumAddr, serviceID string) *Authenticator {
	return &Authenticator{
		jwtSecretKey:    jwtSecretKey,
		node:            node,
		cognitariumAddr: cognitariumAddr,
		serviceID:       serviceID,
	}
}

// Authenticate verifies the provided verifiable credential and issue a related jwt access token if authentication
// succeeds.
func (a *Authenticator) Authenticate(_vc []byte) (string, error) {
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

// Authorize verifies the provided jwt access token
func (a *Authenticator) Authorize(raw string) error {
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

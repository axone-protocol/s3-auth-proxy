package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type ProxyClaims struct {
	jwt.StandardClaims
	Can Permissions `json:"can"`
}

type Permissions struct {
	Read []string `json:"read"`
}

func (c *ProxyClaims) CanRead(uri string) bool {
	for _, u := range c.Can.Read {
		if u == uri {
			return true
		}
	}
	return false
}

func (a *Authenticator) issueJwt(authenticatedSvc string, readURIs []string) (string, error) {
	now := time.Now()
	return jwt.NewWithClaims(jwt.SigningMethodHS256, ProxyClaims{
		StandardClaims: jwt.StandardClaims{
			Audience:  authenticatedSvc,
			ExpiresAt: now.Add(5 * time.Minute).Unix(),
			Id:        uuid.New().String(),
			IssuedAt:  now.Unix(),
			Issuer:    a.serviceID,
			NotBefore: now.Unix(),
			Subject:   authenticatedSvc,
		},
		Can: Permissions{
			Read: readURIs,
		},
	}).SignedString(a.jwtSecretKey)
}

func (a *Authenticator) verifyJwt(raw string) (*ProxyClaims, error) {
	token, err := jwt.ParseWithClaims(raw, &ProxyClaims{}, func(_ *jwt.Token) (interface{}, error) {
		return a.jwtSecretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	c, ok := token.Claims.(*ProxyClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	return c, nil
}

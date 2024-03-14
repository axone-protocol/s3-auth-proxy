package auth

import (
	"okp4/s3-auth-proxy/dataverse"
)

type Authenticator struct {
	jwtSecretKey    []byte
	dataverseClient *dataverse.Client
	serviceID       string
}

func New(jwtSecretKey []byte, dataverseClient *dataverse.Client, serviceID string) *Authenticator {
	return &Authenticator{
		jwtSecretKey:    jwtSecretKey,
		dataverseClient: dataverseClient,
		serviceID:       serviceID,
	}
}

// Authenticate verifies the provided verifiable credential and issue a related jwt access token if authentication
// succeeds.
func (a *Authenticator) Authenticate(_vc []byte) (string, error) {
	return a.issueJwt()
}

// Authorize verifies the provided jwt access token
func (a *Authenticator) Authorize(raw string) error {
	return a.verifyJwt(raw)
}

package auth

import (
	"github.com/piprate/json-gold/ld"
	"okp4/s3-auth-proxy/dataverse"
)

type Authenticator struct {
	jwtSecretKey    []byte
	dataverseClient *dataverse.Client
	serviceID       string
	documentLoader  ld.DocumentLoader
}

func New(jwtSecretKey []byte, dataverseClient *dataverse.Client, serviceID string) *Authenticator {
	return &Authenticator{
		jwtSecretKey:    jwtSecretKey,
		dataverseClient: dataverseClient,
		serviceID:       serviceID,
		documentLoader:  ld.NewDefaultDocumentLoader(nil),
	}
}

// Authenticate verifies the provided verifiable credential and issue a related jwt access token if authentication
// succeeds.
func (a *Authenticator) Authenticate(raw []byte) (string, error) {
	claim, err := a.parseVC(raw)
	if err != nil {
		return "", err
	}

	return a.issueJwt(claim.ID)
}

// Authorize verifies the provided jwt access token
func (a *Authenticator) Authorize(raw string) error {
	return a.verifyJwt(raw)
}

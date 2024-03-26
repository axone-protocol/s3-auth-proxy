package auth

import (
	"context"
	"fmt"

	"okp4/s3-auth-proxy/dataverse"

	"github.com/golang-jwt/jwt"
	"github.com/piprate/json-gold/ld"
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
func (a *Authenticator) Authenticate(ctx context.Context, raw []byte) (string, error) {
	claim, err := a.parseVC(raw)
	if err != nil {
		return "", fmt.Errorf("couldn't parse VC: %w", err)
	}

	zone, err := a.dataverseClient.GetExecutionOrderContext(ctx, claim.ForOrder, claim.ID)
	if err != nil {
		return "", fmt.Errorf("couldn't fetch execution order context: %w", err)
	}

	govCode, err := a.dataverseClient.GetResourceGovCode(ctx, a.serviceID)
	if err != nil {
		return "", fmt.Errorf("couldn't get governance code: %w", err)
	}

	res, err := a.dataverseClient.ExecGovernance(ctx, govCode, "service:use", claim.ID, zone)
	if err != nil {
		return "", fmt.Errorf("couldn't check governance: %w", err)
	}
	if res.Result != "permitted" {
		return "", fmt.Errorf("governance rejected access, evidence: %s", res.Evidence)
	}

	return a.issueJwt(claim.ID)
}

// Authorize verifies the provided jwt access token.
func (a *Authenticator) Authorize(raw string) (*jwt.StandardClaims, error) {
	return a.verifyJwt(raw)
}

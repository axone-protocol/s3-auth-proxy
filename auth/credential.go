package auth

import (
	"fmt"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	"time"
)

type authClaim struct {
	ID        string
	ToService string
	ForOrder  string
}

func (a *Authenticator) parseVC(raw []byte) (*authClaim, error) {
	vc, err := verifiable.ParseCredential(
		raw,
		verifiable.WithJSONLDValidation(),
		verifiable.WithJSONLDDocumentLoader(a.documentLoader),
	)
	if err != nil {
		return nil, err
	}

	if vc.Expired != nil && vc.Expired.After(time.Now()) {
		return nil, fmt.Errorf("verifiable credential expired")
	}

	claim, err := parseAuthClaim(vc)
	if err != nil {
		return nil, err
	}

	if claim.ToService != a.serviceID {
		return nil, fmt.Errorf("auth claim doesn't target us, but service: %s", claim.ToService)
	}

	if vc.Issuer.ID != claim.ID {
		return nil, fmt.Errorf("auth claim subject different from issuer")
	}

	return claim, nil
}

func parseAuthClaim(vc *verifiable.Credential) (*authClaim, error) {
	claims, ok := vc.Subject.([]verifiable.Subject)
	if !ok {
		return nil, fmt.Errorf("malformed vc subject")
	}

	if len(claims) != 1 {
		return nil, fmt.Errorf("expected a single vc claim")
	}

	toService, err := extractCustomStrClaim(&claims[0], "toService")
	if err != nil {
		return nil, err
	}
	forOrder, err := extractCustomStrClaim(&claims[0], "forOrder")
	if err != nil {
		return nil, err
	}

	return &authClaim{
		ID:        claims[0].ID,
		ToService: toService,
		ForOrder:  forOrder,
	}, nil
}

func extractCustomStrClaim(claim *verifiable.Subject, name string) (string, error) {
	field, ok := claim.CustomFields[name]
	if !ok {
		return "", fmt.Errorf("malformed vc claim: '%s' missing", name)
	}

	strField, ok := field.(string)
	if !ok {
		return "", fmt.Errorf("malformed vc claim: expected '%s' to be a string", name)
	}

	return strField, nil
}

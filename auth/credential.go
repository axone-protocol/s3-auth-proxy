package auth

import (
	"fmt"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	"github.com/mitchellh/mapstructure"
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

	var claim authClaim
	if err := mapstructure.Decode(vc.Subject, &claim); err != nil {
		return nil, fmt.Errorf("malformed vc subject: %w", err)
	}

	if claim.ToService != a.serviceID {
		return nil, fmt.Errorf("auth claim doesn't target us, but service: %s", claim.ToService)
	}

	if vc.Issuer.ID != claim.ID {
		return nil, fmt.Errorf("auth claim subject different from issuer")
	}

	return &claim, nil
}

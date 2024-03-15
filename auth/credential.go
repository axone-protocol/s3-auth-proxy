package auth

import (
	"crypto/ecdsa"
	"fmt"
	"time"

	"github.com/btcsuite/btcd/btcec"

	"github.com/hyperledger/aries-framework-go/pkg/doc/jose/jwk/jwksupport"

	"github.com/hyperledger/aries-framework-go/pkg/doc/did"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/verifier"
	"github.com/hyperledger/aries-framework-go/pkg/vdr/fingerprint"

	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/ecdsasecp256k1signature2019"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/ed25519signature2018"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/ed25519signature2020"
	"github.com/hyperledger/aries-framework-go/pkg/vdr"
	"github.com/hyperledger/aries-framework-go/pkg/vdr/key"

	secp "github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
)

type authClaim struct {
	ID        string
	ToService string
	ForOrder  string
}

func (a *Authenticator) parseVC(raw []byte) (*authClaim, error) {
	publicKeyFetcher := verifiable.NewVDRKeyResolver(vdr.New(vdr.WithVDR(key.New()))).PublicKeyFetcher()

	vc, err := verifiable.ParseCredential(
		raw,
		verifiable.WithJSONLDValidation(),
		verifiable.WithPublicKeyFetcher(func(issuerID, keyID string) (*verifier.PublicKey, error) {
			// HACK: as the publicKeyFetcher doesn't manage `EcdsaSecp256k1VerificationKey2019` as verification method
			// we got to manage it ourselves.
			pubKey, err := mayResolveSecp256k1PubKey(issuerID, keyID)
			if err != nil {
				return nil, err
			}

			if pubKey != nil {
				return pubKey, nil
			}

			return publicKeyFetcher(issuerID, keyID)
		}),
		verifiable.WithEmbeddedSignatureSuites(
			ed25519signature2018.New(suite.WithVerifier(ed25519signature2018.NewPublicKeyVerifier())),
			ed25519signature2020.New(suite.WithVerifier(ed25519signature2020.NewPublicKeyVerifier())),
			ecdsasecp256k1signature2019.New(suite.WithVerifier(ecdsasecp256k1signature2019.NewPublicKeyVerifier())),
		),
		verifiable.WithJSONLDDocumentLoader(a.documentLoader),
	)
	if err != nil {
		return nil, err
	}

	if vc.Expired != nil && vc.Expired.After(time.Now()) {
		return nil, fmt.Errorf("verifiable credential expired")
	}

	if len(vc.Proofs) == 0 {
		return nil, fmt.Errorf("missing verifiable credential proof")
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

// Hack helper to resolve a did key as a `EcdsaSecp256k1VerificationKey2019` if it is, as the `PublicKeyFetcher` we use
// doesn't.
func mayResolveSecp256k1PubKey(issuerID, keyID string) (*verifier.PublicKey, error) {
	issuerDid, err := did.Parse(issuerID)
	if err != nil {
		return nil, fmt.Errorf("pub:key vdr Read: failed to parse DID document: %w", err)
	}

	if issuerDid.Method != "key" {
		return nil, fmt.Errorf("vdr Read: invalid did:key method: %s", issuerDid.Method)
	}

	pubKeyBytes, code, err := fingerprint.PubKeyFromFingerprint(issuerDid.MethodSpecificID)
	if err != nil {
		return nil, fmt.Errorf("pub:key vdr Read: failed to get key fingerPrint: %w", err)
	}

	if code == 0xe7 && fmt.Sprintf("#%s", issuerDid.MethodSpecificID) == keyID {
		pubKey, err := secp.ParsePubKey(pubKeyBytes)
		if err != nil {
			return nil, fmt.Errorf("pub:key vdr Read: failed to parse public key: %w", err)
		}
		j, err := jwksupport.JWKFromKey(&ecdsa.PublicKey{
			Curve: btcec.S256(),
			X:     pubKey.X(),
			Y:     pubKey.Y(),
		})
		if err != nil {
			return nil, fmt.Errorf("pub:key vdr Read: error creating JWK: %w", err)
		}

		return &verifier.PublicKey{
			Type: "EcdsaSecp256k1VerificationKey2019",
			JWK:  j,
		}, nil
	}
	return nil, nil
}

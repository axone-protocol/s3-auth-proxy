package auth

import (
	"context"
	"fmt"
	"slices"

	"github.com/valyala/fasthttp"

	"okp4/s3-auth-proxy/dataverse"

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

	execCtx, err := a.dataverseClient.GetExecutionOrderContext(ctx, claim.ForOrder, claim.ID)
	if err != nil {
		return "", fmt.Errorf("couldn't fetch execution order context: %w", err)
	}

	executions := execCtx.ExecutionsInProgress()
	if len(executions) == 0 {
		return "", fmt.Errorf("execution order not in progress")
	}

	var resources []string
	for _, exec := range executions {
		consumed, err := a.dataverseClient.GetExecutionConsumedResources(ctx, claim.ForOrder, exec)
		if err != nil {
			return "", fmt.Errorf("couldn't fetch execution consumed resources: %w", err)
		}

		if !slices.Contains(consumed, a.serviceID) {
			continue
		}
		resources = append(resources, consumed...)
	}

	if len(resources) == 0 {
		return "", fmt.Errorf("not concerned by this execution order")
	}

	res, err := a.execGovernance(ctx, a.serviceID, "service:use", claim.ID, execCtx.Zone)
	if err != nil {
		return "", err
	}
	if res.Result != "permitted" {
		return "", fmt.Errorf("access rejected by governance, evidence: %s", res.Evidence)
	}

	resourcePublications := make([]string, 0)
	for _, r := range resources {
		if r == a.serviceID {
			continue
		}

		uri, err := a.dataverseClient.GetResourcePublication(ctx, r, a.serviceID)
		if err != nil {
			return "", fmt.Errorf("couldn't fetch resource publication: %w", err)
		}
		if uri == nil {
			continue
		}

		res, err := a.execGovernance(ctx, r, "dataset:read", claim.ID, execCtx.Zone)
		if err != nil {
			return "", err
		}
		if res.Result != "permitted" {
			return "", fmt.Errorf("access rejected by governance, evidence: %s", res.Evidence)
		}
		resourcePublications = append(resourcePublications, *uri)
	}

	return a.issueJwt(claim.ID, resourcePublications)
}

// Authorize verifies the provided jwt access token.
func (a *Authenticator) Authorize(token string, uri *fasthttp.URI) (*ProxyClaims, error) {
	claims, err := a.verifyJwt(token)
	if err != nil {
		return nil, err
	}

	if !claims.CanRead(uri.String()) {
		return claims, fmt.Errorf("access to requested resource unauthorized")
	}
	return claims, nil
}

func (a *Authenticator) execGovernance(ctx context.Context, resource, action, subject, zone string) (*dataverse.GovernanceExecAnswer, error) {
	govCode, err := a.dataverseClient.GetResourceGovCode(ctx, resource)
	if err != nil {
		return nil, fmt.Errorf("couldn't fetch governance code: %w", err)
	}

	res, err := a.dataverseClient.ExecGovernance(ctx, govCode, action, subject, zone)
	if err != nil {
		return nil, fmt.Errorf("couldn't exec governance: %w", err)
	}
	return res, nil
}

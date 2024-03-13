package auth

import (
    "fmt"
)

type Authenticator struct {
    node            string
    cognitariumAddr string
    serviceID       string
}

func New(node, cognitariumAddr, serviceID string) *Authenticator {
    return &Authenticator{
        node:            node,
        cognitariumAddr: cognitariumAddr,
        serviceID:       serviceID,
    }
}

// Authenticate verifies the provided verifiable credential and issue a related jwt access token if authentication
// succeeds.
func (a *Authenticator) Authenticate(_vc []byte) (string, error) {
    return "", fmt.Errorf("not implemented")
}

// Authorize verifies the provided jwt access token
func (a *Authenticator) Authorize(jwt string) error {
    return fmt.Errorf("not implemented")
}

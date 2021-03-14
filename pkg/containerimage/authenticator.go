package containerimage

import "github.com/google/go-containerregistry/pkg/authn"

func (i *IdentityAuthenticator) Authorization() (*authn.AuthConfig, error) {
	return &authn.AuthConfig{
		IdentityToken: i.IdentityToken,
	}, nil
}

type IdentityAuthenticator struct {
	IdentityToken string
}

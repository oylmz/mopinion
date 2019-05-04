package mopinion

import (
	"context"
	"encoding/base64"
	"fmt"
)

// TokenInterface has only one method that returns a token by given public and private keys.
type TokenInterface interface {
	Get(ctx context.Context) (*Token, *Response, error)
}

// TokenService implements TokenInterface.
type TokenService struct {
	service
}

// Get returns a token by passing the keys with basic authentication.
func (s *TokenService) Get(ctx context.Context) (*Token, *Response, error) {
	req, err := s.client.NewRequest("GET", "token", nil)
	if err != nil {
		return nil, nil, err
	}

	// Basic authentication https://en.wikipedia.org/wiki/Basic_access_authentication
	authValue := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", s.client.publicKey, s.client.privateKey)))
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", authValue))

	token := new(Token)
	resp, err := s.client.Do(ctx, req, token)
	if err != nil {
		return nil, resp, err
	}
	// We store the token to be used in the future.
	s.client.token = token

	return token, resp, nil
}

package mopinion

import (
	"context"
)

// AccountInterface holds only one method for retrieving the account.
type AccountInterface interface {
	Get(ctx context.Context) (*Account, *Response, error)
}

// AccountService implements AccountInterface.
type AccountService struct {
	service
}

// Get returns the account corresponded to given public key.
func (s *AccountService) Get(ctx context.Context) (*Account, *Response, error) {
	req, err := s.client.NewRequest("GET", "account", nil)
	if err != nil {
		return nil, nil, err
	}
	account := new(Account)
	resp, err := s.client.Do(ctx, req, account)
	if err != nil {
		return nil, resp, err
	}

	return account, resp, nil
}

package mopinion

import (
	"context"
	"testing"
)

func TestToken(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	token, response, err := client.Token.Get(context.Background())
	if err != nil {
		t.Errorf("token API should not return an error: %s", err)
	}

	if response == nil {
		t.Errorf("response can not be nil")
	}

	expectedToken := "token"
	if token.Token != expectedToken {
		t.Errorf("expected token: %v but got: %v", expectedToken, token.Token)
	}
}

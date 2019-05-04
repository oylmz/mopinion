package mopinion

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func setup() (client *Client, mux *http.ServeMux, serverURL string, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"token":"token"}`)
	})

	server := httptest.NewServer(mux)

	client, _ = NewClient(nil, NewBasicCredentialProvider("publickey", "privatekey"))
	url, _ := url.Parse(server.URL + "/")
	client.BaseURL = url

	return client, mux, server.URL, server.Close
}

func TestBasicCredentialProvider(t *testing.T) {
	expectedPublicKey := "testPublicKey"
	expectedPrivateKey := "testPrivateKey"
	basicCredentialProvider := NewBasicCredentialProvider(expectedPublicKey, expectedPrivateKey)
	publicKey, privateKey, err := basicCredentialProvider.Keys()
	if err != nil {
		t.Errorf("basicCredentialProvider should not return an error: %s", err)
	}

	if publicKey != expectedPublicKey {
		t.Errorf("expected public key:%v but got:%v", expectedPublicKey, publicKey)
	}

	if privateKey != expectedPrivateKey {
		t.Errorf("expected private key:%v but got:%v", expectedPrivateKey, privateKey)
	}
}

func TestErrorResponse(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/account", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, `{
			"status": 401,
			"error_code": 18,
			"title": "The credentials you provided are not valid",
			"type": "https://developer.mopinion.com/api/error-codes#invalid-credentials"
		}`)
	})

	// When an incorrect x-auth-token is passed along with the request, the API returns the error above.
	client.privateKey = "incorrectKey"
	client.Token.Get(context.Background())
	account, _, err := client.Account.Get(context.Background())
	if account != nil {
		t.Errorf("account should be nil but got: %v", account)
	}
	if _, ok := err.(*AuthenticationError); !ok {
		t.Errorf("err should be an AuthenticationError but received error: %v", err)
	}
}

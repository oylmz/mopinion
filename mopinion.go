package mopinion

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/google/go-querystring/query"
)

const (
	defaultBaseURL = "https://api.mopinion.com/"
	userAgent      = "mopinion-go-client"
)

const (
	// API error codes.
	// https://developer.mopinion.com/api/error-codes/
	unknown ErrorCode = iota + 1
	server
	publicKeyNotFound
	invalidToken
	invalidJSON
	routeNotFound
	reportIDNotSet
	reportNotFound
	failedToUpdateResource
	notAuthorized
	maxReportsReached
	invalidRequest
	datasetIDNotSet
	datasetNotFound
	maxDatasetsReached
	failedToCreateResource
	noDeploymentCode
	notAuthenticated
)

// ErrorCode represented
type ErrorCode int

// Client represents the Mopinion client.
type Client struct {
	// Client is used to communicate with the mopinion api.
	client *http.Client

	// publicKey is a key used together with the private key to get a token.
	// It has nothing to do with public-key cryptography.
	// It serves like an account id.
	publicKey string

	// privateKey is a key for the public key above.
	privateKey string

	// token is retrieved through token api to be used for
	// generating HMAC signature.
	token *Token

	// BaseURL holds the url for the mopinion api.
	BaseURL *url.URL

	// UserAgent holds the agent name while communicating with mopinion api.
	UserAgent string

	// Services used for talking to different parts of the Mopinion API.
	Token       TokenInterface
	Account     AccountInterface
	Deployments DeploymentsInterface
	Datasets    DatasetsInterface
	Fields      FieldsInterface
	Feedback    FeedbackInterface
	Reports     ReportsInterface
}

type service struct {
	client *Client
}

// CredentialProvider is an interface to be used to retrieve keys to access the Mopinion API.
type CredentialProvider interface {
	// Keys returns public and private keys respectively.
	Keys() (string, string, error)
}

// BasicCredentialProvider implements CredentialProvider.
type BasicCredentialProvider struct {
	publicKey  string
	privateKey string
}

// Keys returns public and private keys respectively.
func (b *BasicCredentialProvider) Keys() (string, string, error) {
	return b.publicKey, b.privateKey, nil
}

// NewBasicCredentialProvider returns a BasicCredentialProvider.
func NewBasicCredentialProvider(publicKey, privateKey string) CredentialProvider {
	return &BasicCredentialProvider{
		publicKey:  publicKey,
		privateKey: privateKey,
	}
}

// NewClient returns a new Mopinion API client.
func NewClient(httpClient *http.Client, credentialProvider CredentialProvider) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	if credentialProvider == nil {
		return nil, fmt.Errorf("credentialProvider cannot be nil")
	}

	baseURL, _ := url.Parse(defaultBaseURL)
	c := &Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent}

	var err error
	if c.publicKey, c.privateKey, err = credentialProvider.Keys(); err != nil {
		return nil, fmt.Errorf("read keys: %s", err)
	}

	service := service{client: c}
	c.Token = &TokenService{service}
	c.Account = &AccountService{service}
	c.Deployments = &DeploymentsService{service}
	c.Datasets = &DatasetsService{service}
	c.Fields = &FieldsService{service}
	c.Feedback = &FeedbackService{service}
	c.Reports = &ReportsService{service}

	return c, nil
}

func (c *Client) makeToken(path string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(string(c.token.Token)))
	mac.Write([]byte(path + "|"))
	mac.Write(body)
	sig := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString([]byte(c.publicKey + ":" + fmt.Sprintf("%x", sig)))
}

// AddAutheticationToken adds an x-auth-token.
func (c *Client) AddAutheticationToken(req *http.Request) error {
	if c.token == nil {
		return fmt.Errorf("token cannot be nil. Get a token first")
	}

	path := req.URL.Path
	// Relative paths may omit leading slash.
	if !strings.HasPrefix(req.URL.Path, "/") {
		path = fmt.Sprintf("/%s", path)
	}
	var body []byte
	if req.Body != nil {
		body, _ = ioutil.ReadAll(req.Body)
	}
	req.Header.Add("x-auth-token", c.makeToken(path, body))
	// Restore the io.ReadCloser to its original state.
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return nil
}

// SetToken sets a token.
func (c *Client) SetToken(token *Token) {
	c.token = token
}

// IsAuthenticationRequired checks if an auth token needs to be passed along with the given relative URL.
func (c *Client) IsAuthenticationRequired(method, urlStr string) bool {
	if method == "GET" && (urlStr == "token" || urlStr == "ping") {
		return false
	}
	return true
}

// NewRequest creates an API request.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a trailing slash, but %q does not", c.BaseURL)
	}
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	if c.IsAuthenticationRequired(method, urlStr) {
		if err := c.AddAutheticationToken(req); err != nil {
			return nil, err
		}
	}

	return req, nil
}

// Do makes a request to mopinion API and returns the response.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	req.WithContext(ctx)
	resp, err := c.client.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, err
	}
	defer resp.Body.Close()
	response := newResponse(resp)
	err = CheckResponse(resp)

	if err != nil {
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err == io.EOF {
				err = nil // ignore EOF errors caused by empty response body
			}
		}
	}
	return response, err
}

// PaginationOptions holds info for pagination.
type PaginationOptions struct {
	Page  int    `url:"page,omitempty"`
	Limit int    `url:"limit,omitempty"`
	Sort  string `url:"sort,omitempty"`
	Order string `url:"order,omitempty"`
}

func addPaginationOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// CheckResponse checks the response if the response is either an error or not.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}
	switch {
	case r.StatusCode == http.StatusUnauthorized && ErrorCode(errorResponse.ErrorCode) == notAuthenticated:
		return (*AuthenticationError)(errorResponse)
	case ErrorCode(errorResponse.ErrorCode) == server:
		return (*ServerError)(errorResponse)
	// TODO: other errors should be added here
	default:
		return errorResponse
	}
}

// AuthenticationError occurs when an invalid token is provided.
type AuthenticationError ErrorResponse

func (r *AuthenticationError) Error() string { return (*ErrorResponse)(r).Error() }

// ServerError is unexpected server side error, coming through Mopinion API.
type ServerError ErrorResponse

func (r *ServerError) Error() string { return "unexpected server side error. Try your request again" }

// ErrorResponse represent the returning error struct.
// Implements the error interface.
type ErrorResponse struct {
	Response  *http.Response
	Status    int    `json:"status"`
	ErrorCode int    `json:"error_code"`
	Title     string `json:"title"`
	Type      string `json:"type"`
}

// Error returns an explainatory message.
func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("method:%v url:%v: http status code:%d error status:%d error code:%d title:%v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Status, r.ErrorCode, r.Title)
}

// Response is wrapper around http.Response.
type Response struct {
	*http.Response
}

func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	return response
}

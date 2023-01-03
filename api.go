package fishfish

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const apiRoot = "https://api.fishfish.gg/v1"

type Category string

const (
	CategorySafe     = "safe"
	CategoryPhishing = "phishing"
	CategoryMalware  = "malware"
)

type Client struct {
	primaryToken string
	sessionToken sessionToken
	permissions  []APIPermission
	apiUrl       string
	httpClient   *http.Client
	// Use defaultAuthType for endpoints where authentication is optional
	defaultAuthType authType
}

func New(primaryToken string, permissions []APIPermission) (*Client, error) {
	client := Client{
		primaryToken: primaryToken,
		permissions:  permissions,
		apiUrl:       apiRoot,
		httpClient:   &http.Client{},
	}

	// Get session token if primaryToken is provided
	if len(primaryToken) > 0 {
		token, err := client.CreateSessionToken()

		if err != nil {
			return nil, fmt.Errorf("failed to create session token: %s", err)
		}

		client.sessionToken = *token
		client.defaultAuthType = authTypeSession
	} else {
		client.defaultAuthType = authTypeNone
	}

	return &client, nil
}

func (c *Client) makeRequest(method, path string, query url.Values, body *bytes.Buffer, authType authType) (*http.Response, error) {

	// Join base and request path
	requestURL, err := url.JoinPath(c.apiUrl, path)

	if err != nil {
		return nil, fmt.Errorf("unable to join path: %s", err)
	}
	// Encode query to string
	queryString := query.Encode()

	fullRequestURL := fmt.Sprintf("%s?%s", requestURL, queryString)

	// Create empty body if null
	if body == nil {
		body = &bytes.Buffer{}
	}

	req, err := http.NewRequest(method, fullRequestURL, body)

	if err != nil {
		return nil, fmt.Errorf("unable to create http request: %s", err)
	}

	// Set auth token based on required authentication type
	switch authType {
	case authTypePrimary:
		req.Header.Set("Authorization", c.primaryToken)
	case authTypeSession:
		req.Header.Set("Authorization", c.sessionToken.Token)
	case authTypeNone:
		// No authorization, do nothing
	}

	// Accept & Send JSON
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := c.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("error sending http request: %s", err)
	}

	// Non-ok status code
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, fmt.Errorf("api returned non-200 status code: %s", res.Status)
	}

	return res, nil
}

// Generic function for marshalling the response into JSON
func readBody[T any](res *http.Response) (*T, error) {
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("could not read response body: %s", err)
	}

	jsonBody := new(T)
	if err = json.Unmarshal(body, jsonBody); err != nil {
		return nil, fmt.Errorf("error converting response body to JSON: %s", err)
	}

	return jsonBody, nil
}

func makeQuery(values map[string]string) url.Values {
	query := url.Values{}
	for k, v := range values {
		query.Add(k, v)
	}

	return query
}

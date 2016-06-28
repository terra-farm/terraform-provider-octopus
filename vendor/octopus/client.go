// Package octopus contains the client for the Octopus Deploy API.
package octopus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

// Client is the client for the Octopus Deploy API.
type Client struct {
	baseAddress string
	stateLock   *sync.Mutex
	httpClient  *http.Client
}

// NewClientWithAPIKey creates a new Octopus Deploy API client using the specified API key.
func NewClientWithAPIKey(serverURL string, apiKey string) (*Client, error) {
	if len(apiKey) == 0 {
		return nil, fmt.Errorf("Must specify a valid Octopus API key.")
	}

	if strings.LastIndex(serverURL, "/") != len(serverURL) {
		serverURL += "/"
	}

	return &Client{
		serverURL,
		&sync.Mutex{},
		&http.Client{
			Transport: &apiKeyAuthenticator{
				apiKey:    apiKey,
				transport: defaultTransport(),
			},
		},
	}, nil
}

// NewClientWithBasicAuth creates a new Octopus Deploy API client using the specified user name and password for HTTP Basic authentication.
func NewClientWithBasicAuth(serverURL string, username string, password string) (*Client, error) {
	if len(username) == 0 {
		return nil, fmt.Errorf("Must specify a valid Octopus user name.")
	}

	return &Client{
		serverURL,
		&sync.Mutex{},
		&http.Client{
			Transport: &usernamePasswordAuthenticator{
				username:  username,
				password:  password,
				transport: defaultTransport(),
			},
		},
	}, nil
}

// Reset clears all cached data from the Client.
func (client *Client) Reset() {
	client.stateLock.Lock()
	defer client.stateLock.Unlock()

	// TODO: Do we actually keep any state in the Octopus client?
}

// Create a request for the octopus API.
func (client *Client) newRequest(relativeURI string, method string, body interface{}) (*http.Request, error) {
	requestURI := fmt.Sprintf("%s/api/%s", client.baseAddress, relativeURI)

	var (
		request    *http.Request
		bodyReader io.Reader
		err        error
	)

	bodyReader, err = newReaderFromJSON(body)
	if err != nil {
		return nil, err
	}

	request, err = http.NewRequest(method, requestURI, bodyReader)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Accept", "application/json")

	if bodyReader != nil {
		request.Header.Set("Content-Type", "application/json")
	}

	return request, nil
}

// executeRequest performs the specified request and returns the entire response body, together with the HTTP status code.
func (client *Client) executeRequest(request *http.Request) (responseBody []byte, statusCode int, err error) {
	if request.Body != nil {
		defer request.Body.Close()
	}

	if os.Getenv("TEST_TF_OCTOPUS_TRACE_HTTP") != "" {
		log.Printf("Invoking %s request for '%s'...", request.Method, request.URL)
	}

	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, 0, err
	}
	defer response.Body.Close()

	statusCode = response.StatusCode

	responseBody, err = ioutil.ReadAll(response.Body)

	if os.Getenv("TEST_TF_OCTOPUS_TRACE_HTTP") != "" {
		log.Printf("Status code: %d, response body: '%s'", statusCode, string(responseBody))
	}

	return
}

// Read an APIErrorResponse (as JSON) from the response body.
func readAPIErrorResponseAsJSON(responseBody []byte, statusCode int) (response *APIErrorResponse, err error) {
	response = &APIErrorResponse{}
	err = json.Unmarshal(responseBody, response)
	if err != nil {
		return
	}

	if len(response.Message) == 0 {
		response.Message = "An unexpected response was received from the Octopus API ('ErrorMessage' field was empty or missing)."
	}

	return
}

// newReaderFromJSON serialises the specified data as JSON and returns an io.Reader over that JSON.
func newReaderFromJSON(data interface{}) (io.Reader, error) {
	if data == nil {
		return nil, nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(jsonData), nil
}

func defaultTransport() *http.Transport {
	return &http.Transport{}
}

// Package octopus contains the client for the Octopus Deploy API.
package octopus

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
)

// Client is the client for the Octopus Deploy API.
type Client struct {
	baseAddress string
	apiKey      string
	username    string
	password    string
	stateLock   *sync.Mutex
	httpClient  *http.Client
}

// NewClientWithAPIKey creates a new Octopus Deploy API client using the specified API key.
func NewClientWithAPIKey(serverURL string, apiKey string) (*Client, error) {
	if len(apiKey) == 0 {
		return nil, fmt.Errorf("Must specify a valid Octopus API key.")
	}

	return &Client{
		serverURL,
		apiKey,
		"",
		"",
		&sync.Mutex{},
		&http.Client{},
	}, nil
}

// NewClientWithBasicAuth creates a new Octopus Deploy API client using the specified user name and password for HTTP Basic authentication.
func NewClientWithBasicAuth(serverURL string, username string, password string) (*Client, error) {
	if len(username) == 0 {
		return nil, fmt.Errorf("Must specify a valid Octopus user name.")
	}

	return &Client{
		serverURL,
		"",
		username,
		password,
		&sync.Mutex{},
		&http.Client{},
	}, nil
}

// Reset clears all cached data from the Client.
func (client *Client) Reset() {
	client.stateLock.Lock()
	defer client.stateLock.Unlock()

	// TODO: Do we actually keep any state in the Octopus client?
}

// executeRequest performs the specified request and returns the entire response body, together with the HTTP status code.
func (client *Client) executeRequest(request *http.Request) (responseBody []byte, statusCode int, err error) {
	if request.Body != nil {
		defer request.Body.Close()
	}

	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, 0, err
	}
	defer response.Body.Close()

	statusCode = response.StatusCode

	responseBody, err = ioutil.ReadAll(response.Body)

	return
}

// Read an APIResponse (as JSON) from the response body.
func readAPIResponseAsJSON(responseBody []byte, statusCode int) (*APIResponse, error) {
	apiResponse := &APIResponse{}
	err := json.Unmarshal(responseBody, apiResponse)
	if err != nil {
		return nil, err
	}

	if len(apiResponse.Message) == 0 {
		apiResponse.Message = "An unexpected response was received from the compute API ('message' field was empty or missing)."
	}

	return apiResponse, nil
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

// newReaderFromXML serialises the specified data as XML and returns an io.Reader over that XML.
func newReaderFromXML(data interface{}) (io.Reader, error) {
	if data == nil {
		return nil, nil
	}

	xmlData, err := xml.Marshal(data)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(xmlData), nil
}

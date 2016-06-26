package octopus

import (
	"net/http"
)

const (
	// HeaderNameOctopusAPIKey is the name of the "X-Octopus-ApiKey" header used for authenticating to Octopus Deploy.
	HeaderNameOctopusAPIKey = "X-Octopus-ApiKey"
)

type apiKeyAuthenticator struct {
	apiKey    string
	transport *http.Transport
}

func (authenticator *apiKeyAuthenticator) RoundTrip(request *http.Request) (*http.Response, error) {
	request.Header[HeaderNameOctopusAPIKey] = []string{authenticator.apiKey}

	return authenticator.transport.RoundTrip(request)
}

type usernamePasswordAuthenticator struct {
	username  string
	password  string
	transport *http.Transport
}

func (authenticator *usernamePasswordAuthenticator) RoundTrip(request *http.Request) (*http.Response, error) {
	request.SetBasicAuth(authenticator.username, authenticator.password)

	return authenticator.transport.RoundTrip(request)
}

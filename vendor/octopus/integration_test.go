package octopus

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
 * Integration test support
 */

type ClientTest struct {
	APIKey      string
	ContentType string
	Invoke      ClientTestRequestInvoker
	Handle      ClientTestRequestHandler
}

func newClientTest() *ClientTest {
	return &ClientTest{
		APIKey:      "my-test-api-key",
		ContentType: "application/json",
	}
}

type ClientTestRequestInvoker func(test *testing.T, client *Client)
type ClientTestRequestHandler func(test *testing.T, request *http.Request) (statusCode int, responseBody string)

func testClientRequest(test *testing.T, clientTest *ClientTest) {
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		expect := expect(test)
		expect.singleHeaderValue(HeaderNameOctopusAPIKey, clientTest.APIKey, request)

		statusCode, response := clientTest.Handle(test, request)

		writer.Header().Set("Content-Type", clientTest.ContentType)
		writer.WriteHeader(statusCode)

		fmt.Fprintln(writer, response)
	}))
	defer testServer.Close()

	client, err := NewClientWithAPIKey(testServer.URL, clientTest.APIKey)
	if err != nil {
		test.Fatal(err)
	}

	clientTest.Invoke(test, client)
}

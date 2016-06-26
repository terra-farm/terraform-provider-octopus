package octopus

import (
	"net/http"
	"reflect"
	"strings"
	"testing"
)

type expectHelper struct {
	test *testing.T
}

func expect(test *testing.T) expectHelper {
	return expectHelper{test}
}

func (expect expectHelper) IsTrue(description string, condition bool) {
	if !condition {
		expect.test.Fatalf("Expression was false: %s", description)
	}
}

func (expect expectHelper) IsFalse(description string, condition bool) {
	if condition {
		expect.test.Fatalf("Expression was true: %s", description)
	}
}

func (expect expectHelper) IsNil(description string, actual interface{}) {
	if !reflect.ValueOf(actual).IsNil() {
		expect.test.Fatalf("%s was not nil.", description)
	}
}

func (expect expectHelper) NotNil(description string, actual interface{}) {
	if reflect.ValueOf(actual).IsNil() {
		expect.test.Fatalf("%s was nil.", description)
	}
}

func (expect expectHelper) EqualsString(description string, expected string, actual string) {
	if actual != expected {
		expect.test.Fatalf("%s was '%s' (expected '%s').", description, actual, expected)
	}
}

func (expect expectHelper) EqualsInt(description string, expected int, actual int) {
	if actual != expected {
		expect.test.Fatalf("%s was %d (expected %d).", description, actual, expected)
	}
}

func (expect expectHelper) singleHeaderValue(headerName string, expected string, request *http.Request) {
	normalisedHeaderName := normaliseHeaderName(headerName)

	headerValues, ok := request.Header[normalisedHeaderName]
	if !ok {
		expect.test.Fatalf("Missing request header '%s'.", headerName)
	}

	if len(headerValues) != 1 {
		expect.test.Fatalf("Request header '%s' has %d values (expected exactly 1).", headerName, len(headerValues))
	}

	actual := headerValues[0]
	expect.EqualsString("Header."+headerName, expected, actual)
}

// Mimic the normalisation of HTTP header names performed by net/http in order to verify them during tests.
func normaliseHeaderName(headerName string) string {
	segments := strings.Split(headerName, "-")
	for index, segment := range segments {
		segments[index] = strings.Title(
			strings.ToLower(segment),
		)
	}

	return strings.Join(segments, "-")
}

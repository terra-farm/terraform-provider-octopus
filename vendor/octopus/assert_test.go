package octopus

import (
	"reflect"
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

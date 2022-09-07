package validator

import (
	"io"
	"log"
	"reflect"
	"testing"
)

func TestNewRegexValidator(t *testing.T) {
	logger := log.New(io.Discard, "", 0)
	validator := NewRegexValidator(logger)
	_, ok := interface{}(validator).(Validator)
	if !ok {
		t.Fatalf("Returned Validator object does not implement Validator interface")
	}
	if validator.l != logger {
		t.Fatalf("Returned Validator object does not have the right logger object: %p != %p", validator.l, logger)
	}
}

func TestValidateLine(t *testing.T) {
	for _, tc := range []struct {
		name string
		in   string
		exp  []string
	}{
		{"empty line", "", nil},
		{"non-matching line", "abcdefghijklmnopqrstuvwxyz", nil},
		{"matching line", "1.1.1.1 - user [1/Jan/1970:00:00:00 +0000] \"GET /master/yoda 1.1\" 200 1234 \"Lynx/2.8\"", []string{"1.1.1.1", "user", "1/Jan/1970:00:00:00 +0000", "GET", "/master/yoda", "1.1", "200", "1234", "Lynx/2.8"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			logger := log.New(io.Discard, "", 0)
			validator := NewRegexValidator(logger)
			tokens := validator.ValidateLine(tc.in)
			if !stringSlicesEqual(tokens, tc.exp) {
				t.Fatalf("Unexpected result of validation: %v != %v", tokens, tc.exp)
			}
		})
	}
}

func TestValidateTokens(t *testing.T) {
	for _, tc := range []struct {
		name string
		in   []string
		expM TokenMap
		expB bool
	}{
		{"empty", nil, TokenMap{}, true},
		{"non-matching 1st", []string{"1.a.1.1"}, TokenMap{}, false},
		{"non-matching 4th", []string{"1.1.1.1", "root", "01/Jan/1970:00:00:00 +0000", "GET01"}, TokenMap{0: "1.1.1.1", 1: "root", 2: "01/Jan/1970:00:00:00 +0000"}, false},
		{"non-matching 8th", []string{"1.1.1.1", "root", "01/Jan/1970:00:00:00 +0000", "GET", "/master/yoda", "HTTP/1.1", "204", "-1123"}, TokenMap{0: "1.1.1.1", 1: "root", 2: "01/Jan/1970:00:00:00 +0000", 3: "GET", 4: "/master/yoda", 5: "HTTP/1.1", 6: "204"}, false},
		{"all-matching", []string{"1.1.1.1", "root", "01/Jan/1970:00:00:00 +0000", "GET", "/master/yoda", "HTTP/1.1", "204", "1123", "some agent"}, TokenMap{0: "1.1.1.1", 1: "root", 2: "01/Jan/1970:00:00:00 +0000", 3: "GET", 4: "/master/yoda", 5: "HTTP/1.1", 6: "204", 7: "1123", 8: "some agent"}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			logger := log.New(io.Discard, "", 0)
			validator := NewRegexValidator(logger)
			tokenMap, validated := validator.ValidateTokens(tc.in)
			if validated != tc.expB {
				t.Fatalf("failed to validate tokens")
			}
			if !reflect.DeepEqual(tokenMap, tc.expM) {
				t.Fatalf("returned token map differs from expected one: %v != %v", tokenMap, tc.expM)
			}
		})
	}
}

func TestGetToken(t *testing.T) {
	for _, tc := range []struct {
		name string
		key  int
		in   TokenMap
		exp  string
	}{
		{"empty map", 0, TokenMap{}, ""},
		{"key not in map", 5, TokenMap{0: "yoda"}, ""},
		{"key in map", 1, TokenMap{0: "yoda", 1: "han solo", 9: "princess leia"}, "han solo"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			logger := log.New(io.Discard, "", 0)
			validator := NewRegexValidator(logger)
			token := validator.GetToken(tc.key, tc.in)
			if token != tc.exp {
				t.Fatalf("retrieved token differs from the expected value: %v != %v", token, tc.exp)
			}
		})
	}
}

func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

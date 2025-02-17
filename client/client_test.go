package client

import (
	"fmt"
	"net/http"
	"testing"
)

func TestExtractCsrfHeaders(t *testing.T) {
	var testCases = []struct {
		headers     http.Header
		expectError bool
	}{
		{
			map[string][]string{
				"X-Harbor-Csrf-Token": {"1234"},
				"Set-Cookie":          {"_gorilla_csrf=5678; expires=in-the-future"},
			},
			false,
		},
		{
			map[string][]string{
				"X-Harbor-Csrf-Token": {"1234"},
				"Set-Cookie":          {"_gorilla_csrf=5678; expires=in-the-future"},
				"Irrelevant-Header":   {"irrelevant"},
			},
			false,
		},
		{
			map[string][]string{
				"X-Harbor-Csrf-Token": {"1234"},
				"Set-Cookie": {
					"abc=9012; irrelevant_cookie=123",
					"_gorilla_csrf=5678; expires=in-the-future",
				},
			},
			false,
		},
		{
			map[string][]string{
				"X-Harbor-Csrf-Token": {"1234", "wrong"},
				"Set-Cookie":          {"_gorilla_csrf=5678; expires=in-the-future"},
			},
			true,
		},
		{
			map[string][]string{
				"X-Harbor-Csrf-Token": {"wrong", "1234"},
				"Set-Cookie":          {"_gorilla_csrf=5678; expires=in-the-future"},
			},
			true,
		},
		{
			map[string][]string{
				"X-Harbor-Csrf-Token": {"1234"},
				"Set-Cookie":          {"_gorilla_csrf=5678"},
			},
			true,
		},
		{
			map[string][]string{
				"X-Harbor-Csrf-Token": {"1234"},
				"Set-Cookie":          {"_gorilla_csrf; expires=in-the-future"},
			},
			true,
		},
		{
			map[string][]string{
				"Set-Cookie": {"_gorilla_csrf=5678; expires=in-the-future"},
			},
			true,
		},
		{
			map[string][]string{
				"X-Harbor-Csrf-Token": {"1234"},
			},
			true,
		},
		{
			map[string][]string{},
			true,
		},
	}

	expected := http.Header{
		"X-Harbor-Csrf-Token": {"1234"},
		"Cookie":              {"_gorilla_csrf=5678"},
	}

	for _, testCase := range testCases {
		testname := fmt.Sprintf("%v,Xerror=%t", testCase.headers, testCase.expectError)
		t.Run(testname, func(t *testing.T) {
			actual, err := extractCsrfHeaders(testCase.headers)
			if testCase.expectError {
				if err == nil {
					t.Fatal("expected error but got nil")
				} else {
					return // Success
				}
			}
			if err != nil {
				t.Fatalf("got unexpected error: %v", err)
			}
			for key, expected_value := range expected {
				actual_value := actual.Get(key)
				if expected_value[0] != actual_value {
					t.Fatalf("%s was %s; want %s", key, expected_value, actual_value)
				}
			}
			for actual_key := range actual {
				if expected.Get(actual_key) == "" {
					t.Fatalf("unexpected key %s", actual_key)
				}
			}
		})
	}
}

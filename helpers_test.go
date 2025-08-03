package main

import (
	"testing"
)

func TestIsUrl(t *testing.T) {
	tests := []struct {
		urlText  string
		expected bool
	}{
		{
			urlText:  "https://example.com",
			expected: true,
		}, {
			urlText:  "http://example.com", // http
			expected: false,
		}, {
			urlText:  "https://example.com:1234", // with port
			expected: false,
		}, {
			urlText:  "https://example.com",
			expected: true,
		}, {
			urlText:  "https://habr.com/ru/companies/ruvds/articles/932220/",
			expected: true,
		}, {
			urlText:  "https://test-example.com",
			expected: true,
		}, {
			urlText:  "https://example.test.com",
			expected: true,
		}, {
			urlText:  "https://test-example.test.com",
			expected: true,
		}, {
			urlText:  "",
			expected: false,
		}, {
			urlText:  "https://chat.deepseek.com/",
			expected: true,
		}, {
			urlText:  "https://example.com:", // with :
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.urlText, func(t *testing.T) {
			result := isUrl(test.urlText)
			if result != test.expected {
				t.Errorf("isUrl(%q) = %v, expected: %v", test.urlText, result, test.expected)
			}
		})
	}
}

func TestExtractDomain(t *testing.T) {
	tests := []struct {
		urlText  string
		expected string
	}{
		{
			urlText:  "https://chat.deepseek.com/",
			expected: "chat.deepseek.com",
		}, {
			urlText: "https://test-example.test.com",
			expected: "test-example.test.com",
		}, {
			urlText: "https://example.com",
			expected: "example.com",
		}, {
			urlText: "https://habr.com/ru/companies/ruvds/articles/932220/",
			expected: "habr.com",
		},
	}

	for _, test := range tests {
		t.Run(test.urlText, func(t *testing.T) {
			result := extractDomain(test.urlText)
			if result != test.expected {
				t.Errorf("extractDomain(%q) = %v, expected: %v", test.urlText, result, test.expected)
			}
		})
	}
}

package utils

import (
	"testing"
)

func TestParseLinksWithValidURL(t *testing.T) {
	data := "http://google.com"
	result := ParseLinks(data)
	if result[0] != "http://google.com" {
		t.Errorf("ParseLink was incorrect, got: %s, want: %s.", result, "http://google.com")
	}
}

func TestParseLinksWithNoURL(t *testing.T) {
	data := "There is no url in this link"
	result := ParseLinks(data)
	if len(result) != 0 {
		t.Errorf("ParseLink was incorrect, got: %d, want: %d.", len(result), 0)
	}
}

func TestParseLinksWithInvalidURL(t *testing.T) {
	data := "www.example.com/main.html"
	result := ParseLinks(data)
	if len(result) != 0 {
		t.Errorf("ParseLink was incorrect, got: %d, want: %d.", len(result), 0)
	}
}

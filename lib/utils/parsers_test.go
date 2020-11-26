package utils

import (
	"os"
	"os/exec"
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

func TestParseIgnoreListPatternWithValidLink(t *testing.T) {
	data := "https://www.example.com/main.html"

	result := ParseIgnoreListPattern(data)

	if len(result) == 0 {
		t.Errorf("ParseIgnoreListPattern was incorrect, got %d, want: %d", len(result), 1)
	}
}

func TestParseIgnoreListPatternWithInvalidLink(t *testing.T) {
	data := "www.example.com/main.html"

	// Mocking an os.Exit(1) scenario
	if os.Getenv("BE_CRASHER") == "1" {
		ParseIgnoreListPattern(data)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestParseIgnoreListPatternWithInvalidLink")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

package utils

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/tonyvugithub/GoURLsCheckerCLI/models"
)

func TestCheckLink(t *testing.T) {
	// Mock a server here
	srv := mockServer(200, "OK")

	// Close the server when done
	defer srv.Close()

	// Set up parameters for CheckLink function
	c := make(chan models.LinkStatus)
	userAgent := "Go-http-client/1.1" //Default User-Agent

	// Waitgroup
	var wg sync.WaitGroup

	wg.Add(1)
	// Make call to the mock server with the server's URL
	go func() {
		defer wg.Done()
		CheckLink(srv.URL, c, userAgent)
	}()

	// Receive the result from the channel
	result := <-c

	wg.Wait()

	resultStatus := result.GetLiveStatus()
	// If the live status is true meaning it returns a 200, else 400 or 404
	if resultStatus != true {
		t.Errorf("CheckLink was incorrect, expect %t, received %t", true, resultStatus)
	}
}

func TestCheckLinkWith404Response(t *testing.T) {
	// Mock a server here
	srv := mockServer(404, "Not Found")

	// Close the server when done
	defer srv.Close()

	// Set up parameters for CheckLink function
	c := make(chan models.LinkStatus)
	userAgent := "Go-http-client/1.1"

	// Waitgroup
	var wg sync.WaitGroup

	wg.Add(1)
	// Make call to the mock server with the server's URL
	go func() {
		CheckLink(srv.URL, c, userAgent)
		defer wg.Done()
	}()

	// Receive the result from the channel
	result := <-c

	wg.Wait()

	resultStatus := result.GetLiveStatus()
	// If the live status is true meaning it returns a 200, else 400 or 404
	if resultStatus != true {
		t.Errorf("CheckLink was incorrect, expect %t, received %t", true, resultStatus)
	}
}

func mockServer(statusCode int, msg string) *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Write([]byte(msg))
	}))
	return srv
}

/* func handler404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 - Not Found"))
}

func handler200(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - OK"))
}
*/

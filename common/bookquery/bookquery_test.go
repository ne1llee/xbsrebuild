package bookquery

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestCachedHttpGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:  "session_id",
			Value: "test-session",
		})
		w.Write([]byte("Hello, world!"))
	}))
	defer server.Close()

	cachedClient, err := NewCachedHTTPClient(true, "cache.json")
	if err != nil {
		t.Fatalf("Error creating CachedHTTPClient: %v", err)
	}

	headers := http.Header{
		"Accept":        []string{"application/json"},
		"Cache-Control": []string{"no-cache"},
	}

	response, err := cachedClient.CachedHttpGet(server.URL, headers, 5*time.Minute)
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}

	expected := "Hello, world!"
	if string(response) != expected {
		t.Errorf("Unexpected response, got: %s, want: %s", response, expected)
	}

	responseFromCache, err := cachedClient.CachedHttpGet(server.URL, headers, 5*time.Minute)
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}

	if string(responseFromCache) != expected {
		t.Errorf("Unexpected response from cache, got: %s, want: %s", responseFromCache, expected)
	}

	// Test session
	serverURL, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("Error parsing server URL: %v", err)
	}

	sessionCookie := cachedClient.Client.Jar.Cookies(serverURL)[0]
	if sessionCookie.Name != "session_id" || sessionCookie.Value != "test-session" {
		t.Errorf("Session cookie not set correctly, got: %+v, want: session_id=test-session", sessionCookie)
	}
}

func TestCachedHttpPost(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("ioutil.ReadAll error: %v", err)
		}
		defer r.Body.Close()
		cookie := &http.Cookie{Name: "test_cookie", Value: "test_value"}
		http.SetCookie(w, cookie)
		fmt.Fprint(w, string(body))
	}))
	defer testServer.Close()

	client, err := NewCachedHTTPClient(true, "")
	if err != nil {
		t.Errorf("NewCachedHTTPClient returned an error: %v", err)
	}

	// Test HTTP POST request
	headers := make(http.Header)
	headers.Set("Test-Header", "test_value")
	data := []byte("test_data")

	response, err := client.CachedHttpPost(testServer.URL, headers, data, time.Second)
	if err != nil {
		t.Errorf("CachedHttpPost returned an error: %v", err)
	}
	if string(response) != "test_data" {
		t.Errorf("CachedHttpPost returned %q, want %q", string(response), "test_data")
	}
}

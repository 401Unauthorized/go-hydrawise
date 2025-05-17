package hydrawise

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Test NewClient with nil and non-nil http.Client
func TestNewClient_DefaultHTTPClient(t *testing.T) {
	c := NewClient("abc", nil)
	if c.HTTPClient != http.DefaultClient {
		t.Error("expected http.DefaultClient")
	}
	if c.APIKey != "abc" {
		t.Error("APIKey not set")
	}
	if !strings.HasSuffix(c.BaseURL, "/api/v1/") {
		t.Error("BaseURL not set")
	}

	customHTTP := &http.Client{}
	c2 := NewClient("key2", customHTTP)
	if c2.HTTPClient != customHTTP {
		t.Error("did not use passed HTTP client")
	}
}

// RoundTrip implementation that always fails
type badRoundTripper struct{}

func (badRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, errors.New("network fail")
}

// doGet error: httpClient.Get fails
func TestClient_doGet_GetFails(t *testing.T) {
	c := NewClient("abc", &http.Client{Transport: badRoundTripper{}})
	var resp interface{}
	err := c.doGet("statusschedule.php", nil, &resp)
	if err == nil || !strings.Contains(err.Error(), "network fail") {
		t.Errorf("expected network fail error, got %v", err)
	}
}

// Simulate a bad Read on Body
type badReadCloser struct{}

func (badReadCloser) Read([]byte) (int, error) { return 0, errors.New("read fail") }
func (badReadCloser) Close() error             { return nil }

// HTTP Transport that always returns a bad Body
type fakeTransport struct{}

func (fakeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body:       badReadCloser{},
	}, nil
}

// doGet error: body read fails
func TestClient_doGet_BadBodyRead(t *testing.T) {
	c := NewClient("abc", &http.Client{Transport: fakeTransport{}})
	var resp interface{}
	err := c.doGet("statusschedule.php", nil, &resp)
	if err == nil || !strings.Contains(err.Error(), "read fail") {
		t.Errorf("expected read fail error, got %v", err)
	}
}

// doGet error: non-200 status code
func TestClient_doGet_Non200(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "forbidden", 403)
	}))
	defer s.Close()
	c := NewClient("abc", s.Client())
	c.BaseURL = s.URL + "/"
	var resp interface{}
	err := c.doGet("", nil, &resp)
	if err == nil || !strings.Contains(err.Error(), "api error") {
		t.Errorf("expected api error, got %v", err)
	}
}

// doGet error: bad JSON
func TestClient_doGet_BadJSON(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		io.WriteString(w, "{bad json")
	}))
	defer s.Close()
	c := NewClient("abc", s.Client())
	c.BaseURL = s.URL + "/"
	var resp interface{}
	err := c.doGet("", nil, &resp)
	if err == nil || !strings.Contains(err.Error(), "invalid character") {
		t.Errorf("expected json unmarshal error, got %v", err)
	}
}

// doGet: check that nil params works
func TestClient_doGet_NilParams(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"msg": "ok"})
	}))
	defer s.Close()
	c := NewClient("abc", s.Client())
	c.BaseURL = s.URL + "/"
	var resp map[string]string
	err := c.doGet("", nil, &resp)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

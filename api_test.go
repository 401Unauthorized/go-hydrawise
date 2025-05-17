package hydrawise

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Helper for mock API server
func mockAPI(t *testing.T, endpoint string, status int, resp interface{}) (*Client, func()) {
	t.Helper()
	handler := http.NewServeMux()
	handler.HandleFunc("/api/v1/"+endpoint, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		switch v := resp.(type) {
		case string:
			w.Write([]byte(v))
		default:
			_ = json.NewEncoder(w).Encode(resp)
		}
	})
	server := httptest.NewServer(handler)
	client := &Client{
		APIKey:     "test-api-key",
		BaseURL:    server.URL + "/api/v1/",
		HTTPClient: server.Client(),
	}
	return client, server.Close
}

// GetStatusSchedule: happy, error (api error), error (bad json)
func TestGetStatusSchedule(t *testing.T) {
	want := &StatusScheduleResponse{
		Time:     1234,
		NextPoll: 60,
		Relays:   []ZoneStatus{{RelayID: 1, Name: "Zone A"}},
	}
	client, closeFn := mockAPI(t, "statusschedule.php", 200, want)
	defer closeFn()

	got, err := client.GetStatusSchedule(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Time != want.Time || len(got.Relays) != 1 {
		t.Errorf("unexpected result: %+v", got)
	}

	// Error: API returns non-200
	clientErr, closeErr := mockAPI(t, "statusschedule.php", 500, `internal error`)
	defer closeErr()
	_, err = clientErr.GetStatusSchedule(nil)
	if err == nil || !strings.Contains(err.Error(), "api error") {
		t.Errorf("expected error, got %v", err)
	}

	// Error: Bad JSON
	clientBadJSON, closeBad := mockAPI(t, "statusschedule.php", 200, `not-json`)
	defer closeBad()
	_, err = clientBadJSON.GetStatusSchedule(nil)
	if err == nil {
		t.Errorf("expected JSON error, got nil")
	}
}

// GetStatusSchedule with controllerID set
func TestGetStatusScheduleWithControllerID(t *testing.T) {
	resp := &StatusScheduleResponse{
		Time:     1,
		NextPoll: 60,
		Relays:   []ZoneStatus{{RelayID: 123, Name: "foo"}},
	}
	client, closeFn := mockAPI(t, "statusschedule.php", 200, resp)
	defer closeFn()
	id := 42
	got, err := client.GetStatusSchedule(&id)
	if err != nil || got.Time != 1 {
		t.Errorf("unexpected: got %v, err %v", got, err)
	}
}

// GetCustomerDetails: happy/error
func TestGetCustomerDetails(t *testing.T) {
	resp := &CustomerDetailsResponse{
		ControllerID:      42,
		CurrentController: "test",
		Controllers:       []ControllerInfo{{Name: "test", SerialNumber: "abc", ControllerID: 42}},
	}
	client, closeFn := mockAPI(t, "customerdetails.php", 200, resp)
	defer closeFn()
	got, err := client.GetCustomerDetails()
	if err != nil || got.ControllerID != 42 {
		t.Errorf("unexpected: %v, %v", got, err)
	}

	// Error: API returns non-200
	clientErr, closeErr := mockAPI(t, "customerdetails.php", 403, `forbidden`)
	defer closeErr()
	_, err = clientErr.GetCustomerDetails()
	if err == nil || !strings.Contains(err.Error(), "api error") {
		t.Errorf("expected error, got %v", err)
	}
}

// RunZone: happy/error
func TestRunZone(t *testing.T) {
	resp := &SetZoneResponse{Message: "Starting zone.", MessageType: "info"}
	client, closeFn := mockAPI(t, "setzone.php", 200, resp)
	defer closeFn()
	got, err := client.RunZone(5, 30)
	if err != nil || got.Message != resp.Message {
		t.Errorf("RunZone failed: got %+v err %v", got, err)
	}

	// Error: API returns error
	clientErr, closeErr := mockAPI(t, "setzone.php", 400, "bad request")
	defer closeErr()
	_, err = clientErr.RunZone(5, 30)
	if err == nil || !strings.Contains(err.Error(), "api error") {
		t.Errorf("expected api error, got %v", err)
	}
}

// StopZone
func TestStopZone(t *testing.T) {
	resp := &SetZoneResponse{Message: "Stopped.", MessageType: "info"}
	client, closeFn := mockAPI(t, "setzone.php", 200, resp)
	defer closeFn()
	got, err := client.StopZone(1)
	if err != nil || got.Message != resp.Message {
		t.Errorf("StopZone failed: got %+v err %v", got, err)
	}
}

// RunAllZones
func TestRunAllZones(t *testing.T) {
	resp := &SetZoneResponse{Message: "Running all zones.", MessageType: "info"}
	client, closeFn := mockAPI(t, "setzone.php", 200, resp)
	defer closeFn()
	got, err := client.RunAllZones(10)
	if err != nil || got.Message != resp.Message {
		t.Errorf("RunAllZones failed: got %+v err %v", got, err)
	}
}

// Nil receiver: for 100% if you have the nil check in code
func TestNilReceiver(t *testing.T) {
	var c *Client
	_, err := c.GetStatusSchedule(nil)
	if err == nil {
		t.Error("expected error from nil receiver")
	}
	_, err = c.GetCustomerDetails()
	if err == nil {
		t.Error("expected error from nil receiver")
	}
	_, err = c.RunZone(1, 30)
	if err == nil {
		t.Error("expected error from nil receiver")
	}
	_, err = c.StopZone(1)
	if err == nil {
		t.Error("expected error from nil receiver")
	}
	_, err = c.RunAllZones(10)
	if err == nil {
		t.Error("expected error from nil receiver")
	}
}

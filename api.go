package hydrawise

import (
	"net/url"
	"strconv"
)

// APIClient defines an interface for the Hydrawise API client.
type APIClient interface {
	GetStatusSchedule(controllerID *int) (*StatusScheduleResponse, error)
	GetCustomerDetails() (*CustomerDetailsResponse, error)
	RunZone(relayID int, seconds int) (*SetZoneResponse, error)
	StopZone(relayID int) (*SetZoneResponse, error)
	RunAllZones(seconds int) (*SetZoneResponse, error)
}

// Ensure Client implements APIClient
var _ APIClient = (*Client)(nil)

// GetStatusSchedule fetches zone schedule/status.
func (c *Client) GetStatusSchedule(controllerID *int) (*StatusScheduleResponse, error) {
	if c == nil {
		return nil, ErrNilClient
	}
	params := url.Values{}
	if controllerID != nil {
		params.Set("controller_id", strconv.Itoa(*controllerID))
	}
	var resp StatusScheduleResponse
	err := c.doGet("statusschedule.php", params, &resp)
	return &resp, err
}

// GetCustomerDetails fetches all controllers for the user.
func (c *Client) GetCustomerDetails() (*CustomerDetailsResponse, error) {
	if c == nil {
		return nil, ErrNilClient
	}
	params := url.Values{}
	var resp CustomerDetailsResponse
	err := c.doGet("customerdetails.php", params, &resp)
	return &resp, err
}

// RunZone starts a zone for given seconds.
func (c *Client) RunZone(relayID int, seconds int) (*SetZoneResponse, error) {
	if c == nil {
		return nil, ErrNilClient
	}
	params := url.Values{}
	params.Set("action", "run")
	params.Set("period_id", "999")
	params.Set("relay_id", strconv.Itoa(relayID))
	params.Set("custom", strconv.Itoa(seconds))
	var resp SetZoneResponse
	err := c.doGet("setzone.php", params, &resp)
	return &resp, err
}

// StopZone stops a given zone.
func (c *Client) StopZone(relayID int) (*SetZoneResponse, error) {
	if c == nil {
		return nil, ErrNilClient
	}
	params := url.Values{}
	params.Set("action", "stop")
	params.Set("relay_id", strconv.Itoa(relayID))
	var resp SetZoneResponse
	err := c.doGet("setzone.php", params, &resp)
	return &resp, err
}

// RunAllZones runs all zones for given seconds.
func (c *Client) RunAllZones(seconds int) (*SetZoneResponse, error) {
	if c == nil {
		return nil, ErrNilClient
	}
	params := url.Values{}
	params.Set("action", "runall")
	params.Set("period_id", "999")
	params.Set("custom", strconv.Itoa(seconds))
	var resp SetZoneResponse
	err := c.doGet("setzone.php", params, &resp)
	return &resp, err
}

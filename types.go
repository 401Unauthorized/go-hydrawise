package hydrawise

import "errors"

// StatusScheduleResponse models statusschedule.php response.
type StatusScheduleResponse struct {
	Time      int64        `json:"time"`
	NextPoll  int          `json:"nextpoll"`
	Message   string       `json:"message"`
	Relays    []ZoneStatus `json:"relays"`
	Sensors   []Sensor     `json:"sensors"`
	SimRelays int          `json:"simRelays"`
	Options   int          `json:"options"`
	StUpdate  int          `json:"stupdate"`
	Expanders []string     `json:"expanders"`
}

// ZoneStatus represents a single zone (relay) status.
type ZoneStatus struct {
	RelayID     int    `json:"relay_id"`
	Time        int    `json:"time"`
	Type        int    `json:"type"`
	Run         int    `json:"run"`
	Relay       int    `json:"relay"`
	Name        string `json:"name"`
	Period      int    `json:"period"`
	TimeStr     string `json:"timestr"`
	Master      *int   `json:"master,omitempty"`
	MasterTimer *int   `json:"master_timer,omitempty"`
}

// Sensor represents a controller sensor.
type Sensor struct {
	Input int `json:"input"`
	Type  int `json:"type"`
}

// CustomerDetailsResponse models customerdetails.php response.
type CustomerDetailsResponse struct {
	ControllerID      int              `json:"controller_id"`
	CustomerID        int              `json:"customer_id"`
	CurrentController string           `json:"current_controller"`
	Controllers       []ControllerInfo `json:"controllers"`
}

// ControllerInfo represents a controller.
type ControllerInfo struct {
	Name         string `json:"name"`
	SerialNumber string `json:"serial_number"`
	ControllerID int    `json:"controller_id"`
}

// SetZoneResponse models setzone.php response.
type SetZoneResponse struct {
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
}

// ErrNilClient is returned when the client is nil.
var ErrNilClient = errors.New("hydrawise client is nil")

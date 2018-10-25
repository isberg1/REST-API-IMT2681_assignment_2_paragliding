package main

// GetIgcinfoAPI used for responding to url: "paragliding/api/"
type GetIgcinfoAPI struct {
	Uptime  string `json:"uptime"`
	Info    string `json:"info"`
	Version string `json:"version"`
}

// InputURL used for receiving IGC url's
type InputURL struct {
	URL string `json:"url"`
}

// ResponsID used as conformation message for responding to client http POST
type ResponsID struct {
	ID string `json:"ID"`
}

// Meta used for storing IGC data from http POST in data structure IgcMap
type Meta struct {
	ID          string `json:"id"`
	TimeStamp   int64  `json:"time_stamp"`
	URL         string `json:"url"`
	HDate       string `json:"h_date"`       //"HDate": <date from File Header, H-record>,
	Pilot       string `json:"pilot"`        //"pilot": <pilot>,
	Glider      string `json:"glider"`       //"glider": <glider>,
	GliderID    string `json:"glider_id"`    //"glider_id": <glider_id>,
	TrackLength int    `json:"track_length"` //"track_length": <calculated total track length>
}

// SimpleMeta is used for responding to http GET requests
type SimpleMeta struct {
	HDate       string `json:"h_date"`       //"HDate": <date from File Header, H-record>,
	Pilot       string `json:"pilot"`        //"pilot": <pilot>,
	Glider      string `json:"glider"`       //"glider": <glider>,
	GliderID    string `json:"glider_id"`    //"glider_id": <glider_id>,
	TrackLength int    `json:"track_length"` //"track_length": <calculated total track length>
	URL         string `json:"url"`
}

// Ticker is used for responding to http GET requests
type Ticker struct {
	TLatest    int64       `json:"t_latest"`
	TStart     int64       `json:"t_start"`
	TStorp     int64       `json:"t_storp"`
	Tracks     []ResponsID `json:"tracks"`
	Processing int64       `json:"processing"`
}

// SimpleWebHookStruct used for registering new webHook subscriptions
type SimpleWebHookStruct struct {
	WebHookURL      string `json:"web_hook_url"`
	MinTriggerValue int    `json:"min_trigger_value"`
}

// WebHookStruct used for storing webhook subscriptions in database
type WebHookStruct struct {
	ID              string `json:"id"`
	TimeStamp       int64  `json:"time_stamp"`
	Counter         int    `json:"counter"`
	WebHookURL      string `json:"web_hook_url"`
	MinTriggerValue int    `json:"min_trigger_value"`
}

// InvokeWebHookStruct used for responding to http GET requests
type InvokeWebHookStruct struct {
	TLatest    int64       `json:"t_latest"`
	Tracks     []ResponsID `json:"tracks"`
	Processing int64       `json:"processing"`
}

type mongoDbStruct struct {
	Host         string
	DatabaseName string
	collection   string
}

/*
//Slack webhook message
type SlackMessage struct {
	Text string `json:"text"`
}
*/

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
	ID string `json:"Id"`
}

// Meta used for storing IGC data from http POST in data structure IgcMap
type Meta struct {
	Id          string `json:"id"`
	TimeStamp   int64  `json:"time_stamp"`
	URL         string `json:"url"`
	HDate       string `json:"h_date"`       //"HDate": <date from File Header, H-record>,
	Pilot       string `json:"pilot"`        //"pilot": <pilot>,
	Glider      string `json:"glider"`       //"glider": <glider>,
	GliderID    string `json:"glider_id"`    //"glider_id": <glider_id>,
	TrackLength int    `json:"track_length"` //"track_length": <calculated total track length>
}

type SimpleMeta struct {
	HDate       string `json:"h_date"`       //"HDate": <date from File Header, H-record>,
	Pilot       string `json:"pilot"`        //"pilot": <pilot>,
	Glider      string `json:"glider"`       //"glider": <glider>,
	GliderID    string `json:"glider_id"`    //"glider_id": <glider_id>,
	TrackLength int    `json:"track_length"` //"track_length": <calculated total track length>
}

// Empty used for returning empty json body when no Id is found in data structure IgcMap
type Empty struct {
	_ string `json:""`
}

type Ticker struct {
	TLatest    int64       `json:"t_latest"`
	TStart     int64       `json:"t_start"`
	TStorp     int64       `json:"t_storp"`
	Tracks     []ResponsID `json:"tracks"`
	Processing int64       `json:"processing"`
}

type SimpleWebHookStruct struct {
	WebHookURL      string `json:"web_hook_url"`
	MinTriggerValue int    `json:"min_trigger_value"`
}
type WebHookStruct struct {
	Id              string `json:"id"`
	TimeStamp       int64  `json:"time_stamp"`
	Counter         int    `json:"counter"`
	WebHookURL      string `json:"web_hook_url"`
	MinTriggerValue int    `json:"min_trigger_value"`
}
type InvokeWebHookStruct struct {
	TLatest    int64       `json:"t_latest"`
	Tracks     []ResponsID `json:"tracks"`
	Processing int64       `json:"processing"`
}

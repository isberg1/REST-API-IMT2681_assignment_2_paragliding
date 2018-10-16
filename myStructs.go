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
	ID string `json:"id"`
}

// Meta used for storing IGC data from http POST in data structure IgcMap
type Meta struct {
	HDate       string `json:"h_date"`       //"HDate": <date from File Header, H-record>,
	Pilot       string `json:"pilot"`        //"pilot": <pilot>,
	Glider      string `json:"glider"`       //"glider": <glider>,
	GliderID    string `json:"glider_id"`    //"glider_id": <glider_id>,
	TrackLength int    `json:"track_length"` //"track_length": <calculated total track length>
}

// Empty used for returning empty json body when no ID is found in data structure IgcMap
type Empty struct {
	_ string `json:""`
}

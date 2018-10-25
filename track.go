package main

type Track struct {
	HDate string `json:"H_date"`
	Pilot string `json:"pilot"`
	Glider string `json:"glider"`
	GliderID string `json:"glider_id"`
	TrackLength float64 `json:"track_length"`
	TrackSrcUrl string `json:"track_src_url"`
}

type TrackData struct {
	Id int `json:"id"`
	Track Track `json:"track"`
}

type TrackDB struct {
	Tracks []TrackData `json:"track"`
}
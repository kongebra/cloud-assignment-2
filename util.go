package main

import (
	"fmt"
	"os"
	"time"
)

import (
	"github.com/marni/goigc"
	"math"
)

/**
	Code found at: https://gist.github.com/cdipaolo/d3f8db3848278b49db68
 */
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

/**
	Code found at : https://gist.github.com/cdipaolo/d3f8db3848278b49db68
 */
func Distance(lat1, lon1, lat2, lon2 float64) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}

// Create Track from URL and igc.Track
func CreateTrackFromIGC(url string, t igc.Track) Track {
	// Get all points of the igc.Track
	points := t.Points

	// Get first latitude & longitude
	lat := points[0].Lat
	lng := points[0].Lng

	// Declare length variable
	var length float64

	// Loop through points
	for i, p := range points {
		// Check that are not at the first points
		if i != 0 {
			// Add the distance from points to the length
			length += Distance(float64(lat), float64(lng), float64(p.Lat), float64(p.Lng))

			// Get the next latitude & longitude
			lat = p.Lat
			lng = p.Lng
		}
	}

	// Returns a Track
	return Track{
		HDate: t.Date.String(),
		Pilot: t.Pilot,
		Glider: t.GliderType,
		GliderID: t.GliderID,
		TrackLength: length,
		TrackSrcUrl: url,
		Timestamp: time.Now().Unix(),
	}
}

// Converts seconds to ISO-8601 (duration)
func ConvertSecondsToISO8601(seconds int) string {
	sec := seconds % 60
	min := seconds / 60
	hour := min / 60
	days := hour / 24
	month := days / 30
	year := month / 12

	min %= 60
	hour %= 24
	days %= 30
	month %= 12

	result := "P"

	if year > 0 {
		result += fmt.Sprintf("%dY", year)
	}

	if month > 0 {
		result += fmt.Sprintf("%dM", month)
	}

	if days > 0 {
		result += fmt.Sprintf("%dD", days)
	}

	if hour > 0 || min > 0 || sec > 0 {
		result += "T"

		if hour > 0 {
			result += fmt.Sprintf("%dH", hour)
		}

		if min > 0 {
			result += fmt.Sprintf("%dM", min)
		}

		if sec > 0 {
			result += fmt.Sprintf("%dS", sec)
		}
	}


	return result
}

// Gets port from the environment
func GetPort() string {
	var port = os.Getenv("PORT")

	// Check if port is blank (localhost)
	if port == "" {
		port = "4747"
	}

	return ":" + port
}
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"
)

var apiKey = os.Getenv("NASA_API_KEY")

type NeoResponse struct {
	NearEarthObjects map[string][]NeoObject `json:"near_earth_objects"`
}

type NeoObject struct {
	ID                   string `json:"id"`
	PotentiallyHazardous bool   `json:"is_potentially_hazardous_asteroid"`
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// Get start and end date for a given calendar week
func getDateIntervalByWeek(calendarWeek int) (startDate time.Time, endDate time.Time) {
	year, _ := time.Now().ISOWeek()

	// January 4th is always in the first week of the year
	jan4 := time.Date(year, 1, 4, 0, 0, 0, 0, time.UTC)

	// Normalize the weekday to 1 (Monday) - 7 (Sunday)
	weekday := int(jan4.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	// Search for the Monday of the first week
	mondayWeek1 := jan4.AddDate(0, 0, -(weekday - 1))

	startDate = mondayWeek1.AddDate(0, 0, (calendarWeek-1)*7)
	endDate = startDate.AddDate(0, 0, 6)

	return startDate, endDate
}

// Get astroid IDs for a given calendar week
func getIDs(w http.ResponseWriter, r *http.Request) {
	cwStr := r.URL.Path[len("/neo/week/"):]

	cw, err := strconv.Atoi(cwStr)
	if err != nil {
		http.Error(w, "invalid calendar week", http.StatusBadRequest)
		return
	}

	year := time.Now().Year()
	maxWeek := maxISOWeek(year)

	if cw < 1 || cw > maxWeek {
		http.Error(w, fmt.Sprintf("calendar week must be 1-%d", maxWeek), http.StatusBadRequest)
		return
	}

	startDate, endDate := getDateIntervalByWeek(cw)

	url := fmt.Sprintf(
		"https://api.nasa.gov/neo/rest/v1/feed?start_date=%s&end_date=%s&api_key=%s",
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"),
		apiKey,
	)

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "NASA request failed", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "NASA API error", http.StatusBadGateway)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "read failed", http.StatusInternalServerError)
		return
	}

	var neoResp NeoResponse
	if err := json.Unmarshal(body, &neoResp); err != nil {
		http.Error(w, "decode failed", http.StatusInternalServerError)
		return
	}

	// FILTER: hazardous optional
	filter := r.URL.Query().Get("filter")
	onlyHazardous := filter == "hazardous"

	var ids []string

	for _, objects := range neoResp.NearEarthObjects {
		for _, obj := range objects {
			if onlyHazardous && !obj.PotentiallyHazardous {
				continue
			}
			ids = append(ids, obj.ID)
		}
	}

	sort.Strings(ids)

	if ids == nil {
		ids = []string{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ids)
}

func maxISOWeek(year int) int {
	t := time.Date(year, 12, 28, 0, 0, 0, 0, time.UTC)
	_, week := t.ISOWeek()
	return week
}

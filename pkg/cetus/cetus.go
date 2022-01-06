package cetus

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

// WorldStateJSON is used to extract json infos
type WorldStateJSON struct {
	SyndicateMissions []SyndicateMissionJSON `json:"SyndicateMissions"`
}

// SyndicateMissionJSON is used to extract json infos
type SyndicateMissionJSON struct {
	Tag       string   `json:"Tag"`
	StartDate DateJSON `json:"Activation"`
	EndDate   DateJSON `json:"Expiry"`
}

// DateJSON is used to extract json infos
type DateJSON struct {
	Date NumberJSON `json:"$date"`
}

// NumberJSON is used to extract json infos
type NumberJSON struct {
	Timestamp string `json:"$numberLong"`
}

// Time is the time struct for Cetus
type Time struct {
	DayStart   time.Time
	NightStart time.Time
	NightEnd   time.Time
}

// World contains all the retrieved infos
var World = &Time{}

// Represents the duration of the time cycles
const (
	Day     = 100 * time.Minute
	Night   = 50 * time.Minute
	FullDay = 150 * time.Minute
)

// Start populates the variables for the time
func Start() {
	populateTime()

	tk := time.NewTicker(24 * time.Hour)

	go func() {
		for {
			select {
			case <-tk.C:
				populateTime()
			}
		}
	}()
}

// PopulateTime is used to retrieve Cetus infos from warframe servers
func populateTime() {
	var worldState = &WorldStateJSON{}

	resp, err := http.Get("https://content.warframe.com/dynamic/worldState.php")
	if err != nil {
		log.Fatalln("Error connecting to website", err)
	}

	err = json.NewDecoder(resp.Body).Decode(&worldState)

	for _, v := range worldState.SyndicateMissions {
		if v.Tag == "CetusSyndicate" {
			start, _ := strconv.ParseInt(v.StartDate.Date.Timestamp, 10, 64)
			end, _ := strconv.ParseInt(v.EndDate.Date.Timestamp, 10, 64)

			World.DayStart = time.UnixMilli(start)
			World.NightStart = time.UnixMilli(start).Add(time.Minute * 100)
			World.NightEnd = time.UnixMilli(end)
		}
	}
}

// WorldTime return the string time before the night appear
func WorldTime() float64 {
	durationPassed := time.Duration(time.Since(World.DayStart)/FullDay) * FullDay

	if math.Mod(time.Since(World.DayStart).Seconds(), 150*60) < 100*60 {
		// Day
		return time.Until(World.NightStart.Add(durationPassed)).Truncate(1 * time.Minute).Minutes()
	}

	return 0
}

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

// Cetus contains all the retrieved infos
var Cetus = &Time{}

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

			Cetus.DayStart = time.UnixMilli(start)
			Cetus.NightStart = time.UnixMilli(start).Add(time.Minute * 100)
			Cetus.NightEnd = time.UnixMilli(end)
		}
	}
}

// WorldTime return the string time before the night appear
func WorldTime() float64 {
	daysPassed := time.Duration(time.Since(Cetus.DayStart).Seconds() / (150 * 60) * float64(time.Second))

	if math.Mod(time.Since(Cetus.DayStart).Seconds(), 150*60) < 100*60 {
		// Day
		return time.Until(Cetus.NightStart.Add(daysPassed)).Truncate(1 * time.Minute).Minutes()
	}

	return 0
}

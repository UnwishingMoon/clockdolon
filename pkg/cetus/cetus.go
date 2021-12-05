package cetus

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

type WorldStateJson struct {
	SyndicateMissions []SyndicateMissionJson `json:"SyndicateMissions"`
}

type SyndicateMissionJson struct {
	Tag       string   `json:"Tag"`
	StartDate DateJson `json:"Activation"`
	EndDate   DateJson `json:"Expiry"`
}

type DateJson struct {
	Date NumberJson `json:"$date"`
}

type NumberJson struct {
	Timestamp string `json:"$numberLong"`
}

type CetusTime struct {
	DayStart   time.Time
	NightStart time.Time
	NightEnd   time.Time
}

var Cetus = &CetusTime{}

func PopulateCetusTime() {
	var worldState = &WorldStateJson{}

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

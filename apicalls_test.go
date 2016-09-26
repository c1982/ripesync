package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

const (
	TEST_URL_SPRINT        = "https://stat.ripe.net/data/announced-prefixes/data.json?resource=%s"
	TEST_URL_WITH_RESOURCE = "https://stat.ripe.net/data/announced-prefixes/data.json?resource=AS43260"
	TEXT_JSON_PREFIX_DATA  = `{
    "status": "ok", 
    "server_id": "stat-app8", 
    "status_code": 200, 
    "version": "1.2", 
    "cached": false, 
    "see_also": [], 
    "time": "2016-09-24T20:19:51.089556", 
    "messages": [
        [
            "info", 
            "Results exclude routes with very low visibility (less than 3 RIS full-feed peers seeing)."
        ]
    ], 
    "data_call_status": "supported - connecting to ursa", 
    "process_time": 409, 
    "build_version": "2016.9.9.138", 
    "query_id": "3f19a486-8294-11e6-81db-0050568835e6", 
    "data": {
        "resource": "43260", 
        "prefixes": [ 
            {
                "timelines": [
                    {
                        "endtime": "2016-09-24T16:00:00", 
                        "starttime": "2016-09-10T16:00:00"
                    }
                ], 
                "prefix": "185.119.80.0/24"
            }, 
            {
                "timelines": [
                    {
                        "endtime": "2016-09-24T16:00:00", 
                        "starttime": "2016-09-10T16:00:00"
                    }
                ], 
                "prefix": "185.86.152.0/24"
            }, 
            {
                "timelines": [
                    {
                        "endtime": "2016-09-24T16:00:00", 
                        "starttime": "2016-09-10T16:00:00"
                    }
                ], 
                "prefix": "185.28.62.0/24"
            }, 
            {
                "timelines": [
                    {
                        "endtime": "2016-09-24T16:00:00", 
                        "starttime": "2016-09-10T16:00:00"
                    }
                ], 
                "prefix": "185.93.52.0/24"
            }, 
            {
                "timelines": [
                    {
                        "endtime": "2016-09-24T16:00:00", 
                        "starttime": "2016-09-10T16:00:00"
                    }
                ], 
                "prefix": "185.85.236.0/24"
            }, 
            {
                "timelines": [
                    {
                        "endtime": "2016-09-24T16:00:00", 
                        "starttime": "2016-09-10T16:00:00"
                    }
                ], 
                "prefix": "31.210.159.0/24"
            }, 
            {
                "timelines": [
                    {
                        "endtime": "2016-09-24T16:00:00", 
                        "starttime": "2016-09-10T16:00:00"
                    }
                ], 
                "prefix": "185.86.14.0/24"
            }, 
            {
                "timelines": [
                    {
                        "endtime": "2016-09-24T16:00:00", 
                        "starttime": "2016-09-10T16:00:00"
                    }
                ], 
                "prefix": "185.86.13.0/24"
            }, 
            {
                "timelines": [
                    {
                        "endtime": "2016-09-24T16:00:00", 
                        "starttime": "2016-09-10T16:00:00"
                    }
                ], 
                "prefix": "185.122.202.0/24"
            }, 
            {
                "timelines": [
                    {
                        "endtime": "2016-09-24T16:00:00", 
                        "starttime": "2016-09-10T16:00:00"
                    }
                ], 
                "prefix": "185.136.207.0/24"
            }
        ], 
        "query_starttime": "2016-09-10T16:00:00", 
        "latest_time": "2016-09-24T16:00:00", 
        "query_endtime": "2016-09-24T16:00:00", 
        "earliest_time": "2000-08-01T00:00:00"
    }
}`
)

/*
func TestProdGetJSonData(t *testing.T) {

	text, err := getJsonData(TEST_URL)

	if err != nil {
		t.Errorf("Cannot get data from Ripe API", text)
	}

	if text == "" {
		t.Errorf("Json data is empty", text)
	}

	//t.Log(text)
}
*/

func TestAnouncmentUnMarshalling(t *testing.T) {

	anon := Announcement{}
	err := json.Unmarshal([]byte(TEXT_JSON_PREFIX_DATA), &anon)

	if err != nil {
		t.Errorf("Unmarshalling Error:", err)
	}

	prefix_lenght := len(anon.Data.Prefixes)
	excpected_lenght := 10
	if prefix_lenght != excpected_lenght {
		t.Error("Prefix ")
	}

	/*
		for _, prf := range anon.Data.Prefixes {
			t.Log(prf.Name)
		}
	*/
}

func TestGenerateRangeArray(t *testing.T) {
	list := getScanRange("AS43260")

	for i := 0; i < len(list); i++ {
		fmt.Println(list[i])
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	ANNOUNCED_URL = "https://stat.ripe.net/data/announced-prefixes/data.json?resource=%s"
	RESOURCES_URL = "https://stat.ripe.net/data/country-resource-list/data.json?resource=%s"
)

func getJsonData(uri string) (string, error) {
	resp, err := http.Get(uri)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	jsondata, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(jsondata), err
}

func getPrefixes(autnum string) (Announcement, error) {

	anon := Announcement{}
	uri := fmt.Sprintf(ANNOUNCED_URL, autnum)

	log.Println(uri)

	jsonData, err := getJsonData(uri)

	if err != nil {
		return anon, err
	}

	err = json.Unmarshal([]byte(jsonData), &anon)

	if err != nil {
		return anon, err
	}

	return anon, err
}

func getAsNumbers(country string) (Announcement, error) {
	anon := Announcement{}
	uri := fmt.Sprintf(RESOURCES_URL, country)

	log.Println(uri)

	jsonData, err := getJsonData(uri)

	if err != nil {
		return anon, err
	}

	err = json.Unmarshal([]byte(jsonData), &anon)

	if err != nil {
		return anon, err
	}

	return anon, err
}

func getRangeArrayForConfig(autnum string) ([]string, error) {

	ipcidr := []string{}
	anon, err := getPrefixes(autnum)

	if err != nil {
		return ipcidr, err
	}

	for _, prf := range anon.Data.Prefixes {
		if isCidrIpV4(prf.Name) {
			ipcidr = append(ipcidr, fmt.Sprintf("range=%s", prf.Name))
		}
	}

	return ipcidr, err

}

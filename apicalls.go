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

func getPrefixes(asn string) (ipv4Prefixes []Prefix, ipv6Prefixes []Prefix, err error) {

	anon := Announcement{}
	uri := fmt.Sprintf(ANNOUNCED_URL, asn)

	log.Println(uri)

	jsonData, err := getJsonData(uri)

	if err != nil {
		return nil, nil, err
	}

	err = json.Unmarshal([]byte(jsonData), &anon)

	if err != nil {
		return nil, nil, err
	}

	for _, p := range anon.Data.Prefixes {
		if isCidrIpV4(p.Name) {
			ipv4Prefixes = append(ipv4Prefixes, p)
		} else {
			ipv6Prefixes = append(ipv6Prefixes, p)
		}
	}

	return ipv4Prefixes, ipv6Prefixes, err
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

func GenerateRangeForConfigFile(autnum string, prefixes []Prefix) []string {

	ipcidr := []string{}

	for _, prf := range prefixes {
		ipcidr = append(ipcidr, fmt.Sprintf("range=%s", prf.Name))
	}

	return ipcidr

}

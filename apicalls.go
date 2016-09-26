package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	ANNOUNCED_URL = "https://stat.ripe.net/data/announced-prefixes/data.json?resource=%s"
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

	fmt.Println(uri)

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

func getScanRange(autnum string) []string {

	ipcidr := []string{}
	anon, err := getPrefixes(autnum)

	if err != nil {
		return ipcidr
	}

	for _, prf := range anon.Data.Prefixes {
		if isCidrIpV4(prf.Name) {
			ipcidr = append(ipcidr, "range="+prf.Name)
		}
	}

	return ipcidr

}
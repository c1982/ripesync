package main

import (
	"encoding/json"
	_ "io"
	"os"
	"strings"
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

	TEXT_JSON_RESOURCE_DATA = `{
    "status": "ok", 
    "server_id": "stat-app9", 
    "status_code": 200, 
    "version": "0.2", 
    "cached": false, 
    "see_also": [], 
    "time": "2016-09-26T14:57:07.247951", 
    "messages": [], 
    "data_call_status": "supported - connecting to ursa", 
    "process_time": 918, 
    "build_version": "2016.9.26.140", 
    "query_id": "7ddfc572-83f9-11e6-b268-0050568836e2", 
    "data": {
        "query_time": "2016-09-25T00:00:00", 
        "resources": {
            "ipv6": [
                "2001:678:1a4::/48", 
                "2001:67c:464::/48", 
                "2001:67c:4f4::/48", 
                "2001:67c:5e0::/48", 
                "2001:67c:68c::/48", 
                "2001:67c:6c0::/48", 
                "2001:67c:1154::/48", 
                "2001:67c:11b8::/48", 
                "2001:67c:11ec::/48", 
                "2001:67c:12a4::/48"                
            ], 
            "asn": [
                "1885", 
                "2592", 
                "2600", 
                "2872", 
                "3188", 
                "5422", 
                "5458", 
                "5474", 
                "6707", 
                "6755"                              
            ], 
            "ipv4": [
                "5.2.80.0/21", 
                "5.11.128.0/17", 
                "5.23.120.0/21", 
                "5.24.0.0/14", 
                "5.44.80.0/20", 
                "5.44.144.0/20", 
                "5.46.0.0/15", 
                "5.63.32.0/19", 
                "5.104.0.0/20", 
                "5.159.248.0/21"                          
            ]
        }
    }
}`

	TEXT_JSON_REGISTRY_BROWSER = `{
    "status": "ok", 
    "server_id": "stat-app8", 
    "status_code": 200, 
    "version": "0.3", 
    "cached": false, 
    "see_also": [], 
    "time": "2016-09-27T20:09:23.736710", 
    "messages": [], 
    "data_call_status": "supported", 
    "process_time": 148, 
    "build_version": "2016.9.26.140", 
    "query_id": "489056de-84ee-11e6-aea0-0050568835e6", 
    "data": {
        "num_versions": 38, 
        "resource": "AS42926", 
        "backward_refs": [], 
        "database": "RIPE", 
        "versions": [
            [
                "2007-05-10T11:58:49", 
                "2009-10-16T11:42:48"
            ], 
            [
                "2016-08-15T06:24:55", 
                "2016-09-27T00:00:00"
            ]
        ], 
        "suggestions": [], 
        "complete_history": false, 
        "objects": [
            {
                "fields": [
                    {
                        "references": [], 
                        "value": "AS42926", 
                        "key": "aut-num"
                    }, 
                    {
                        "references": [], 
                        "value": "RADORE", 
                        "key": "as-name"
                    }, 
                    {
                        "references": [
                            "organisation"
                        ], 
                        "value": "ORG-RHTH1-RIPE", 
                        "key": "org"
                    }, 
                    {
                        "references": [
                            "person", 
                            "role"
                        ], 
                        "value": "RLA11-RIPE", 
                        "key": "admin-c"
                    }, 
                    {
                        "references": [
                            "person", 
                            "role"
                        ], 
                        "value": "RLA11-RIPE", 
                        "key": "tech-c"
                    }, 
                    {
                        "references": [
                            "mntner"
                        ], 
                        "value": "MNT-AS42926", 
                        "key": "mnt-by"
                    }, 
                    {
                        "references": [
                            "mntner"
                        ], 
                        "value": "RIPE-NCC-END-MNT", 
                        "key": "mnt-by"
                    }, 
                    {
                        "references": [
                            "mntner"
                        ], 
                        "value": "MNT-AS42926", 
                        "key": "mnt-routes"
                    }
                ], 
                "type": "aut-num", 
                "primary": {
                    "value": "AS42926", 
                    "key": "aut-num"
                }
            }
        ], 
        "time": "2016-09-27T20:09:23", 
        "forward_refs": [
            {
                "relationships": [
                    "org"
                ], 
                "fields": [
                    {
                        "references": [], 
                        "value": "ORG-RHTH1-RIPE", 
                        "key": "organisation"
                    }, 
                    {
                        "references": [], 
                        "value": "Radore Veri Merkezi Hizmetleri A.S.", 
                        "key": "org-name"
                    }, 
                    {
                        "references": [], 
                        "value": "LIR", 
                        "key": "org-type"
                    }, 
                    {
                        "references": [
                            "person", 
                            "role"
                        ], 
                        "value": "YD868-RIPE", 
                        "key": "admin-c"
                    }, 
                    {
                        "references": [
                            "person", 
                            "role"
                        ], 
                        "value": "ZKA2-RIPE", 
                        "key": "admin-c"
                    }, 
                    {
                        "references": [], 
                        "value": "RARA7-RIPE", 
                        "key": "abuse-c"
                    }, 
                    {
                        "references": [
                            "mntner"
                        ], 
                        "value": "RIPE-NCC-HM-MNT", 
                        "key": "mnt-ref"
                    }, 
                    {
                        "references": [
                            "mntner"
                        ], 
                        "value": "AS42926-MNT", 
                        "key": "mnt-ref"
                    }, 
                    {
                        "references": [
                            "mntner"
                        ], 
                        "value": "RIPE-NCC-HM-MNT", 
                        "key": "mnt-by"
                    }, 
                    {
                        "references": [
                            "mntner"
                        ], 
                        "value": "AS42926-MNT", 
                        "key": "mnt-by"
                    }
                ], 
                "type": "organisation", 
                "primary": {
                    "value": "ORG-RHTH1-RIPE", 
                    "key": "organisation"
                }
            }, 
            {
                "relationships": [
                    "admin-c", 
                    "tech-c"
                ], 
                "fields": [
                    {
                        "references": [], 
                        "value": "RADORE LIR", 
                        "key": "role"
                    }, 
                    {
                        "references": [
                            "organisation"
                        ], 
                        "value": "ORG-RHTH1-RIPE", 
                        "key": "org"
                    }, 
                    {
                        "references": [
                            "person", 
                            "role"
                        ], 
                        "value": "RNOC6-RIPE", 
                        "key": "admin-c"
                    }, 
                    {
                        "references": [
                            "person", 
                            "role"
                        ], 
                        "value": "RNOC6-RIPE", 
                        "key": "tech-c"
                    }, 
                    {
                        "references": [], 
                        "value": "RLA11-RIPE", 
                        "key": "nic-hdl"
                    }, 
                    {
                        "references": [], 
                        "value": "abuse@radore.com", 
                        "key": "abuse-mailbox"
                    }, 
                    {
                        "references": [
                            "mntner"
                        ], 
                        "value": "AS42926-MNT", 
                        "key": "mnt-by"
                    }
                ], 
                "type": "role", 
                "primary": {
                    "value": "RLA11-RIPE", 
                    "key": "nic-hdl"
                }
            }, 
            {
                "relationships": [
                    "mnt-by", 
                    "mnt-routes"
                ], 
                "fields": [
                    {
                        "references": [], 
                        "value": "MNT-AS42926", 
                        "key": "mntner"
                    }, 
                    {
                        "references": [], 
                        "value": "AS42926", 
                        "key": "descr"
                    }, 
                    {
                        "references": [
                            "person", 
                            "role"
                        ], 
                        "value": "RLA11-RIPE", 
                        "key": "admin-c"
                    }, 
                    {
                        "references": [
                            "mntner"
                        ], 
                        "value": "MNT-AS42926", 
                        "key": "mnt-by"
                    }
                ], 
                "type": "mntner", 
                "primary": {
                    "value": "MNT-AS42926", 
                    "key": "mntner"
                }
            }, 
            {
                "relationships": [
                    "mnt-by"
                ], 
                "fields": [
                    {
                        "references": [], 
                        "value": "RIPE-NCC-END-MNT", 
                        "key": "mntner"
                    }, 
                    {
                        "references": [], 
                        "value": "RIPE-NCC IPRA Maintainer for all direct resources", 
                        "key": "descr"
                    }, 
                    {
                        "references": [
                            "person", 
                            "role"
                        ], 
                        "value": "CREW-RIPE", 
                        "key": "admin-c"
                    }, 
                    {
                        "references": [
                            "person", 
                            "role"
                        ], 
                        "value": "CREW-RIPE", 
                        "key": "tech-c"
                    }, 
                    {
                        "references": [
                            "organisation"
                        ], 
                        "value": "ORG-NCC1-RIPE", 
                        "key": "org"
                    }, 
                    {
                        "references": [
                            "mntner"
                        ], 
                        "value": "RIPE-NCC-END-MNT", 
                        "key": "mnt-by"
                    }
                ], 
                "type": "mntner", 
                "primary": {
                    "value": "RIPE-NCC-END-MNT", 
                    "key": "mntner"
                }
            }
        ], 
        "now": "2016-09-27T20:09:23", 
        "type": "aut-num", 
        "truncated": false
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

func TestGenerateRangeArrayForMassscan(t *testing.T) {
	list := getScanRange("AS43260")

	for i := 0; i < len(list); i++ {
		fmt.Println(list[i])
	}
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

func TestResourcesUnMarshalling(t *testing.T) {

	anon := Announcement{}
	err := json.Unmarshal([]byte(TEXT_JSON_RESOURCE_DATA), &anon)

	if err != nil {
		t.Errorf("Unmarshalling Error:", err)
	}

	excpected_lenght := 10
	asnumber_lenght := len(anon.Data.Resources.ASNumbers)

	if asnumber_lenght != excpected_lenght {
		t.Error("as number array error")
	}

	ipv4_lenght := len(anon.Data.Resources.IPv4)
	if ipv4_lenght != excpected_lenght {
		t.Error("ipv4 array error")
	}

	ipv6_lenght := len(anon.Data.Resources.IPv6)
	if ipv6_lenght != excpected_lenght {
		t.Error("ipv6 array error")
	}
}

func TestRegistryBrowserUnmarshalling(t *testing.T) {
	anon := Announcement{}
	err := json.Unmarshal([]byte(TEXT_JSON_REGISTRY_BROWSER), &anon)

	if err != nil {
		t.Errorf("Unmarshalling Error:", err)
	}

	field := anon.Data.Objects[0].Fields[1]

	if field.ObjKey != "as-name" {
		t.Error("Invalid object key:", field.ObjKey)
	}

	if field.ObjValue != "RADORE" {
		t.Error("Invalid object value:", field.ObjValue)
	}

	FrwRef := anon.Data.ForwardRefs[0].Fields[1]
	t.Log(FrwRef.ObjValue)
	if FrwRef.ObjKey != "org-name" {
		t.Error("Invalid object key:", FrwRef.ObjKey)
	}
}

func TestReadTemplateFile(t *testing.T) {
	fileTxt, err := readTemplateFile()

	if err != nil {
		t.Error(err)
	}

	if fileTxt == "" {
		t.Error("File cannot be read.")
	}
}

func TestGenerateConfigFile(t *testing.T) {

	anon := Announcement{}
	err := json.Unmarshal([]byte(TEXT_JSON_PREFIX_DATA), &anon)

	if err != nil {
		t.Errorf("Unmarshalling Error:", err)
	}

	fileTxt, _ := readTemplateFile()

	configFile, err := GenerateConfig("AS43260", anon.Data.Prefixes)

	t.Logf("Config file name:", configFile)

	if err != nil {
		t.Errorf("Generate config error:", err)
	}

	if fileTxt == "" {
		t.Error("Template file is empty.")
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Fatalf("File not found:", configFile, err)
	}

	configTxt, err := readFile(configFile)

	if err != nil {
		t.Error(err)
	}

	if configTxt == "" {
		t.Errorf("ASN config is empty", configFile)
	}

	if strings.HasSuffix(configTxt, "#IP RANGE") {
		t.Error("Cannot append prefix list in config file.")
	}

	err = deleteFile(configFile)

	if err != nil {
		t.Errorf("ASN config delete error:", err, configFile)
	}

}

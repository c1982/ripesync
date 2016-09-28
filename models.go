package main

import (
	"time"
)

type Netnum struct {
	Inetnum  string `json:"inetnum"`
	Netname  string `json:"netname"`
	Desc     string `json:"desc"`
	Country  string `json:"country"`
	Admin    string `json:"admin-c"`
	Tech     string `json:"tech-c"`
	Notify   string `json:"notify"`
	MntBy    string `json:"mnt-by"`
	NetStart string
	NetEnd   string
}

type Domain struct {
	Host       string `json:"domain"`
	Desc       string `json:"desc"`
	NameServer string `json:"nserver"`
	MntBy      string `json:"mnt-by"`
}

type Host struct {
	Ip        string
	AutNum    string
	AsName    string
	Org       string
	MntRoutes string
	Prfx      string
}

type HostActivity struct {
	Asn   string
	Ip    string `json:"ip"`
	Date  time.Time
	Ports []Prt `json:"ports"`
}

type Prt struct {
	Port    int    `json:"port"`
	Proto   string `json:"proto"`
	Status  string `json:"open"`
	Reason  string `json:"reason"`
	Ttl     int    `json:"ttl"`
	Service Srv    `json:"service"`
}

type Srv struct {
	Name   string `json:"name"`
	Banner string `json:"banner"`
}

type MntNer struct {
	Mntner   string `json:"mntner"`
	AdminC   string `json:"admin-c"`
	Auth     string `json:"auth"`
	MntBy    string `json:"mnt-by"`
	Created  string `json:"created"`
	Modified string `json:"last-modified"`
}

type Routing struct {
	Route  string `json:"route"`
	Descr  string `json:"descr"`
	Origin string `json:"origin"`
	MntBy  string `json:"mnt-by"`
}

//aut-num
type AsNumber struct {
	AutNum    string
	AsName    string
	Org       string
	MntRoutes string
}

//as-set
type AsProfile struct {
	AsSet string
	Descr string
}

type Announcement struct {
	Status         string        `json:"status"`
	Server         string        `json:"server_id"`
	StatusCode     int           `json:"status_code"`
	Version        string        `json:"version"`
	Cached         bool          `json:"cached"`
	Time           string        `json:"time"`
	DataCallStatus string        `json:"data_call_status"`
	ProcessTime    int           `json:"process_time"`
	BuildVersion   string        `json:"build_version"`
	QueryId        string        `json:"query_id"`
	Data           AnnouncedData `json:"data"`
}

func (r *Announcement) GetObjValue(keyName string) string {

	result := ""
	if len(r.Data.Objects) > 0 {
		for _, f := range r.Data.Objects[0].Fields {
			if f.ObjKey == keyName {
				result = f.ObjValue
				break
			}
		}
	}

	return result
}

func (r *Announcement) GetFrwValue(keyName string) string {
	result := ""
	if len(r.Data.ForwardRefs) > 0 {
		for _, f := range r.Data.ForwardRefs[0].Fields {
			if f.ObjKey == keyName {
				result = f.ObjValue
				break
			}
		}
	}

	return result
}

type AnnouncedData struct {
	Resource       string       `json:"resource"`
	Prefixes       []Prefix     `json:"prefixes"`
	Resources      ResourceData `json:"resources"`
	QueryStarttime string       `json:"query_starttime"`
	LatestTime     string       `json:"latest_time"`
	QueryEndtime   string       `json:"query_endtime"`
	EarliestTime   string       `json:"earliest_time"`
	Objects        []DataObject `json:"objects"`
	ForwardRefs    []DataObject `json:"forward_refs"`
}

type DataObject struct {
	Fields []ObjectField `json:"fields"`
}

type ObjectField struct {
	ObjValue string `json:"value"`
	ObjKey   string `json:"key"`
}

type ResourceData struct {
	IPv6      []string `json:"ipv6"`
	ASNumbers []string `json:"asn"`
	IPv4      []string `json:"ipv4"`
}

type Prefix struct {
	TimeLines []TimeLine `json:"timelines"`
	Name      string     `json:"prefix"`
}

type TimeLine struct {
	EndTime   string `json:"endtime"`
	StartTime string `json:"starttime"`
}

type Summary struct {
	Date time.Time

	Asn       string
	AsnName   string
	Org       string
	OrgName   string
	OrgType   string
	MntRoutes string

	TotalPrefix int
	Ipv4Prefix  int
	Ipv6Prefix  int
	TotalIpv4   int
	ActiveIp4   float64

	Windows  float64
	Linux    float64
	RouterOs float64

	Ftp        float64
	Ssh        float64
	Telnet     float64
	Smtp       float64
	SmtpAuth   float64
	Dns        float64
	Http       float64
	Pop3       float64
	Imap       float64
	Snmp       float64
	Rdp        float64
	Sip        float64
	PowerShell float64
	WebDeploy  float64

	VestaCp      float64
	DirectAdmin  float64
	Plesk        float64
	WebsitePanel float64
	MaestroPanel float64
	CPanel       float64
	CPanelSSL    float64
	CPanelWHM    float64
	CPanelWHMSSL float64
	Ajenti       float64
	Webmin       float64
	HstCntr      float64

	MsSQL      float64
	MySQL      float64
	MongoDB    float64
	PostgreSQL float64
	Redis      float64
}

type TotalAsnStat struct {
	Country    string
	AsnCount   int
	Ipv4Prefix int
	Ipv6Prefix int
	Date       time.Time
}

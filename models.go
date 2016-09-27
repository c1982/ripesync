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

type AnnouncedData struct {
	Resource       string       `json:"resource"`
	Prefixes       []Prefix     `json:"prefixes"`
	Resources      ResourceData `json:"resources"`
	QueryStarttime string       `json:"query_starttime"`
	LatestTime     string       `json:"latest_time"`
	QueryEndtime   string       `json:"query_endtime"`
	EarliestTime   string       `json:"earliest_time"`
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
	Date        time.Time
	Asn         string
	AsnName     string
	Description string

	TotalPrefix int
	Ipv4Prefix  int
	Ipv6Prefix  int
	TotalIpv4   int
	ActiveIp4   int

	Windows  int
	Linux    int
	RouterOs int

	Ftp        int
	Ssh        int
	Telnet     int
	Smtp       int
	Dns        int
	Http       int
	Pop3       int
	Imap       int
	Snmp       int
	Rdp        int
	Sip        int
	PowerShell int
	WebDeploy  int

	WestaCp      int
	DirectAdmin  int
	Plesk        int
	WebsitePanel int
	MaestroPanel int
	CPanel       int
	CPanelSSL    int
	CPanelWHM    int
	CPanelWHMSSL int
	Ajenti       int
	Webmin       int
	HstCntr      int

	MsSQL      int
	MySQL      int
	MongoDB    int
	PostgreSQL int
	Redis      int
}

type TotalAsnStat struct {
	Country    string
	AsnCount   int
	Ipv4Prefix int
	Ipv6Prefix int
	Date       time.Time
}

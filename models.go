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
	Resource       string   `json:"resource"`
	Prefixes       []Prefix `json:"prefixes"`
	QueryStarttime string   `json:"query_starttime"`
	LatestTime     string   `json:"latest_time"`
	QueryEndtime   string   `json:"query_endtime"`
	EarliestTime   string   `json:"earliest_time"`
}

type Prefix struct {
	TimeLines []TimeLine `json:"timelines"`
	Name      string     `json:"prefix"`
}

type TimeLine struct {
	EndTime   string `json:"endtime"`
	StartTime string `json:"starttime"`
}

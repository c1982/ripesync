package main

import (
	"time"
)

type Company struct {
	Inetnum string `json:"inetnum"`
	Netname string `json:"netname"`
	Desc    string `json:"desc"`
	Country string `json:"country"`
	MntBy   string `json:"mnt-by"`
}

type Domain struct {
	Host       string `json:"domain"`
	Desc       string `json:"desc"`
	NameServer string `json:"nserver"`
	MntBy      string `json:"mnt-by"`
}

type Host struct {
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

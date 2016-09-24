package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
)

const (
	rpsl_line_pattern    = `(.+):\W+(.+)`
	ripe_db_inetnum_file = "/home/ripe.db.inetnum"
	rip_db_name          = "ipstat"
)

func main() {

	session, err := mgo.Dial("localhost")

	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	log.Println("Begin")
	SyncRipeDb(*session, ripe_db_inetnum_file, "inetnum:", "inetnum", InsertInetNum)

	log.Println("End")
}

func SyncScanFile(session mgo.Session) {

	file, err := os.OpenFile(os.Args[1], os.O_RDONLY, os.ModeExclusive)
	defer file.Close()

	if err != nil {
		fmt.Println(err)
	}
	c := session.DB(rip_db_name).C("hostactivity")

	scan := bufio.NewScanner(file)

	for scan.Scan() {

		i := HostActivity{}
		line := strings.TrimSpace(scan.Text())

		if strings.HasSuffix(line, ",") {
			line = strings.TrimRight(line, ",")
		}

		byteLine := []byte(line)
		err = json.Unmarshal(byteLine, &i)

		if err != nil {
			fmt.Println(err.Error())
		} else {
			i.Date = time.Now()

			if len(i.Ports) > 0 {
				c.Insert(&i)
			}
		}
	}
}

func InsertInetNum(c mgo.Collection, aggrate string) {

	var d = Netnum{}
	d.Inetnum = parseRPSLValue(aggrate, "inetnum", "inetnum")
	d.Netname = parseRPSLValue(aggrate, "inetnum", "netname")
	d.Desc = parseRPSLValue(aggrate, "inetnum", "descr")
	d.Country = parseRPSLValue(aggrate, "inetnum", "country")
	d.Admin = parseRPSLValue(aggrate, "inetnum", "admin-c")
	d.Tech = parseRPSLValue(aggrate, "inetnum", "tech-c")
	d.Notify = parseRPSLValue(aggrate, "inetnum", "notify")
	d.MntBy = parseRPSLValue(aggrate, "inetnum", "mnt-by")

	if d.Inetnum != "" {

		d.Country = strings.ToUpper(d.Country)

		if d.Country != "TR" {
			return
		}

		log.Println(d.Netname + ":" + d.Inetnum)

		d.NetStart = GetIpBlock(d.Inetnum, 0)
		d.NetEnd = GetIpBlock(d.Inetnum, 1)

		err := c.Insert(&d)

		if err != nil {
			fmt.Println(err)
		}

		n := c.Database.C("hosts")

		expand := ExpandRage(d.NetStart, d.NetEnd)

		for i := 0; i <= len(expand)-1; i++ {

			var h = Host{}
			h.Ip = expand[i]
			h.Admin = d.Admin
			h.Country = d.Country
			h.Inetnum = d.Inetnum
			h.MntBy = d.MntBy
			h.Netname = d.Netname
			h.Notify = d.Notify

			err = n.Insert(&h)

			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func InsertDomain(c mgo.Collection, aggrate string) {

	var d = Domain{}
	d.Host = parseRPSLValue(aggrate, "domain", "domain")
	d.Desc = parseRPSLValue(aggrate, "domain", "descr")
	d.NameServer = parseRPSLValue(aggrate, "domain", "nserver")
	d.MntBy = parseRPSLValue(aggrate, "domain", "mnt-by")

	if d.Host != "" {
		err := c.Insert(&d)

		if err != nil {
			fmt.Println(err)
		}
	}
}

func InsertMntner(c mgo.Collection, aggrate string) {

	var d = MntNer{}
	d.Mntner = parseRPSLValue(aggrate, "mntner", "mntner")
	d.AdminC = parseRPSLValue(aggrate, "mntner", "admin-c")
	d.Auth = parseRPSLValue(aggrate, "mntner", "auth")
	d.MntBy = parseRPSLValue(aggrate, "mntner", "mnt-by")
	d.Created = parseRPSLValue(aggrate, "mntner", "created")
	d.Modified = parseRPSLValue(aggrate, "mntner", "last-modified")

	if d.Mntner != "" {
		err := c.Insert(&d)

		if err != nil {
			fmt.Println(err)
		}
	}
}

func InsertRoute(c mgo.Collection, aggrate string) {

	var d = Routing{}
	d.Route = parseRPSLValue(aggrate, "route", "route")
	d.Descr = parseRPSLValue(aggrate, "route", "descr")
	d.Origin = parseRPSLValue(aggrate, "route", "origin")
	d.MntBy = parseRPSLValue(aggrate, "route", "mnt-by")

	if d.Route != "" {

		expand, err := ExpandRoute(d.Route)

		if err != nil {
			fmt.Println(err)
		}

		for i := 0; i <= len(expand); i++ {

		}

		err = c.Insert(&d)

		if err != nil {
			fmt.Println(err)
		}
	}
}

func SyncRipeDb(session mgo.Session, dbfile string, split string, collectionName string, Insert func(c mgo.Collection, aggrate string)) {

	aggrate := false
	aggrate_string := ""

	c := session.DB(rip_db_name).C(collectionName)

	file, err := os.OpenFile(dbfile, os.O_RDONLY, 0)

	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, split) {
			aggrate = true
		}

		if aggrate {
			if line != "" {
				aggrate_string += fmt.Sprintf("%s\n", line)
			}
		}

		if line == "" {

			Insert(*c, aggrate_string)

			aggrate = false
			aggrate_string = ""
		}
	}
}

func parseRPSLValue(whoisText string, class string, section string) string {

	var sectionValue = ""
	var hasIn = false

	sc := bufio.NewScanner(strings.NewReader(whoisText))
	sc.Split(bufio.ScanLines)

	for sc.Scan() {
		var line = sc.Text()

		if strings.HasPrefix(line, class) {
			if hasIn == false {
				hasIn = true
			}
		}

		if hasIn {
			if strings.HasPrefix(line, section) {
				if sectionValue != "" {
					sectionValue += " "
				}

				sectionValue += parseRPSLine(line)
			}
		}
	}

	return sectionValue
}

func parseRPSLine(whoisLine string) string {

	rx, _ := regexp.Compile(rpsl_line_pattern)
	s := rx.FindAllStringSubmatch(whoisLine, -1)

	if len(s) >= 1 {
		return s[0][2]
	}

	return ""
}

func GetIpBlock(inetnum string, i int) string {
	return strings.TrimSpace(strings.Split(inetnum, "-")[i])
}

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
	rpsl_line_pattern   = `(.+):\W+(.+)`
	ripe_db_domain_file = "O:\\ripe.db.domain\\ripe.db.domain"
	ripe_db_mntner_file = "O:\\ripe.db.domain\\ripe.db.mntner"
	ripe_db_route_file  = "O:\\ripe.db.domain\\ripe.db.route"
	rip_db_name         = "ipstat"
)

func main() {

	session, err := mgo.Dial("localhost")

	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	log.Println("Begin")

	//SyncRipeDb(*session, ripe_db_domain_file, "domain:", "domain", InsertDomain)
	//SyncRipeDb(*session, ripe_db_mntner_file, "mntner:", "mntner", InsertMntner)
	SyncRipeDb(*session, ripe_db_mntner_file, "route:", "route", InsertMntner)

	log.Println("End")
}

func SyncScanFile() {

	file, err := os.OpenFile(os.Args[1], os.O_RDONLY, os.ModeExclusive)
	defer file.Close()

	if err != nil {
		fmt.Println(err)
	}

	scan := bufio.NewScanner(file)

	for scan.Scan() {

		i := Host{}
		line := scan.Text()

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
				AddRow(i)
			}
		}
	}
}

func AddRow(h Host) {
	//Save to Mongo
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

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
)

const (
	rpsl_line_pattern = `(.+):\W+(.+)`
	rip_db_name       = "ipstat"
	template_file     = "main.config.template"
)

func main() {

	country := flag.String("country", "TR", "Country Code: TR, IT, DE etc.")
	scanAsn := flag.Bool("scan", true, "Scan Country As Numbers")
	syncdb := flag.Bool("sync", false, "Check Ripe Database file and save database.")
	autnum := flag.String("autnum", "ripe.db.aut-num", "Ripe aut-num db file.")

	session, err := mgo.Dial("localhost")

	if err != nil {
		log.Println("Can't connect mongodb.")
		panic(err)
	}

	session.DB(rip_db_name).Login("admin", "")
	session.SetMode(mgo.Monotonic, true)
	defer session.Close()

	if *scanAsn {
		if *country == "" {
			log.Println("Required country parameter.")
			return
		}

		log.Printf("Country:", *country)

		log.Println("Begin as scanning.")
		ScanAsNumbers(*country, *session)
		log.Println("End as scanning.")
	}

	if *syncdb {
		if *autnum == "" {
			log.Println("Required autnum parameter.")
			return
		}

		log.Println("Begin autnum sync.")
		SyncRipeDb(*session, *autnum, "aut-num:", "autnums", InsertAutNum)
		log.Println("End autnum sync.")
	}
}

func ScanAsNumbers(country string, session mgo.Session) {
	anon, err := getAsNumbers(country)

	if err != nil {
		panic(err)
	}

	log.Println("AS number count:", len(anon.Data.Resources.ASNumbers))

	for _, v := range anon.Data.Resources.ASNumbers {

		asn := fmt.Sprintf("AS%s", v)
		log.Println("Scanning: ", asn)

		//asnFile := fmt.Sprintf("%s.config", asn)
		scanoutput := fmt.Sprintf("%s.json", asn)
		cnfFile, err := GenerateConfig(asn)

		if err != nil {
			//deleteFile(asnFile)
			log.Println("Config file cannot genereted.", err)
			continue
		}

		//masscan
		err = executeScan(cnfFile, scanoutput)
		if err != nil {
			log.Println("Scan error: ", err)
			//deleteFile(asnFile)
			continue
		}

		//save mongodb
		SyncScanFile(scanoutput, session, asn)

		//deleteFile(asnFile)
		//deleteFile(scanoutput)
	}
}

func SyncScanFile(scanFile string, session mgo.Session, asnum string) {

	file, err := os.OpenFile(scanFile, os.O_RDONLY, 0)
	defer file.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	c := session.DB(rip_db_name).C("hostactivity")

	scan := bufio.NewScanner(file)
	scan.Split(bufio.ScanLines)

	for scan.Scan() {

		line := strings.TrimSpace(scan.Text())

		if len(line) > 2048 {
			continue
		}

		if line == "{finished: 1}" {
			continue
		}

		if strings.HasSuffix(line, ",") {
			line = strings.TrimRight(line, ",")
		}

		h := HostActivity{}
		byteLine := []byte(line)

		err = json.Unmarshal(byteLine, &h)

		if err != nil {
			fmt.Println(err)
			continue
		}
		h.Asn = asnum
		h.Date = time.Now()

		if len(h.Ports) > 0 {
			err = c.Insert(&h)

			if err != nil {
				fmt.Println(err)
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
		/*
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
		*/
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

		_, err := ExpandRoute(d.Route)

		if err != nil {
			fmt.Println(err)
		}

		err = c.Insert(&d)

		if err != nil {
			fmt.Println(err)
		}
	}
}

func InsertAutNum(c mgo.Collection, aggrate string) {

	var d = AsNumber{}
	d.AutNum = parseRPSLValue(aggrate, "aut-num", "aut-num") //AS Number
	d.AsName = parseRPSLValue(aggrate, "aut-num", "as-name")
	d.MntRoutes = parseRPSLValue(aggrate, "aut-num", "mnt-routes")
	d.Org = parseRPSLValue(aggrate, "aut-num", "org")

	err := c.Insert(&d)

	if err != nil {
		fmt.Println(err)
		return
	}

	//Get Prefixes from Ripe
	anon, err := getPrefixes(d.AutNum)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(len(anon.Data.Prefixes))

	if len(anon.Data.Prefixes) <= 0 {
		return
	}

	n := c.Database.C("hosts")

	for _, prf := range anon.Data.Prefixes {

		if isCidrIpV4(prf.Name) {
			fmt.Println(prf.Name)

			//Expand prefix ip address
			prefixIpList, err := ExpandRoute(prf.Name)

			if err != nil {
				fmt.Println(err)
			}

			for i := 0; i < len(prefixIpList); i++ {

				//Insert IP
				h := Host{}
				h.Ip = prefixIpList[i]
				h.AsName = d.AsName
				h.AutNum = d.AutNum
				h.MntRoutes = d.MntRoutes
				h.Org = d.Org
				h.Prfx = prf.Name

				err = n.Insert(&h)

				if err != nil {
					fmt.Println(err)
				}

			}
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

func GenerateConfig(asn string) (string, error) {

	fileName := fmt.Sprintf("%s.config", asn)

	tmplTxt, err := readTemplateFile()

	if err != nil {
		log.Println("Read template error.")
		return "", err
	}

	prfx, err := getRangeArrayForConfig(asn)

	if err != nil {
		return "", err
	}

	log.Println("Prefix count is:", len(prfx))

	if len(prfx) <= 0 {
		msg := fmt.Sprintf("Prefixes is empty this ASN: %s", asn)
		return "", fmt.Errorf(msg)
	}

	rangeLines := strings.Join(prfx, "\n")
	tmplTxt = fmt.Sprintf("%s\n%s", tmplTxt, rangeLines)

	err = CreateFile(fileName)

	if err != nil {
		log.Println("File create error:", fileName)
		return "", err
	}

	err = WriteAllText(fileName, tmplTxt)

	if err != nil {
		log.Println("File write error:", fileName)
		return "", err
	}

	return fileName, err
}

func readTemplateFile() (string, error) {
	return readFile(template_file)
}

func readFile(filePath string) (string, error) {
	file, err := ioutil.ReadFile(filePath)

	if err != nil {
		return "", err
	}

	return string(file), err
}

func deleteFile(filePath string) error {
	return os.Remove(filePath)
}

func executeScan(cnfFile string, outputfile string) error {

	configFile := fmt.Sprintf("/home/masscan/bin/%s", cnfFile)
	outputFileName := fmt.Sprintf("/home/masscan/bin/%s", outputfile)

	log.Printf("Waiting for scan. File: ", cnfFile)
	out, err := exec.Command("/home/masscan/bin/masscan", "--conf", configFile, "-oJ", outputFileName).Output()
	log.Printf("Output: %s", out)

	return err
}

func CreateFile(filePath string) error {
	cnf, err := os.Create(filePath)
	defer cnf.Close()

	return err
}

func WriteAllText(filePath string, text string) error {
	var file, err = os.OpenFile(filePath, os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString(text)

	if err != nil {
		return err
	}

	err = file.Sync()

	return err
}

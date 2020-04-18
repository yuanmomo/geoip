// GeoIP generator
//
// Before running this file, the GeoIP database must be downloaded and present.
// To download GeoIP database: https://dev.maxmind.com/geoip/geoip2/geolite2/
// Inside you will find block files for IPv4 and IPv6 and country code mapping.
package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"geoip/list"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/protobuf/proto"
	"v2ray.com/core/app/router"
	"v2ray.com/core/common"
	"v2ray.com/core/infra/conf"
)

var (
	countryCodeFile = flag.String("country", "", "Path to the country code file")
	ipv4File        = flag.String("ipv4", "", "Path to the IPv4 block file")
	ipv6File        = flag.String("ipv6", "", "Path to the IPv6 block file")
)

func getCountryCodeMap() (map[string]string, error) {
	countryCodeReader, err := os.Open(*countryCodeFile)
	if err != nil {
		return nil, err
	}
	defer countryCodeReader.Close()

	m := make(map[string]string)
	reader := csv.NewReader(countryCodeReader)
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	for _, line := range lines[1:] {
		id := line[0]
		countryCode := line[4]
		if len(countryCode) == 0 {
			continue
		}
		m[id] = strings.ToUpper(countryCode)
	}
	return m, nil
}

func getCidrPerCountry(file string, m map[string]string, list map[string][]*router.CIDR) error {
	fileReader, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fileReader.Close()

	reader := csv.NewReader(fileReader)
	lines, err := reader.ReadAll()
	if err != nil {
		return err
	}
	for _, line := range lines[1:] {
		cidrStr := line[0]
		countryId := line[1]
		if countryCode, found := m[countryId]; found {
			cidr, err := conf.ParseIP(cidrStr)
			if err != nil {
				return err
			}
			cidrs := append(list[countryCode], cidr)
			list[countryCode] = cidrs
		}
	}
	return nil
}

func main() {
	flag.Parse()


	ccMap, err := getCountryCodeMap()
	if err != nil {
		fmt.Println("Error reading country code map:", err)
		return
	}

	cidrList := make(map[string][]*router.CIDR)
	if err := getCidrPerCountry(*ipv4File, ccMap, cidrList); err != nil {
		fmt.Println("Error loading IPv4 file:", err)
		return
	}
	if err := getCidrPerCountry(*ipv6File, ccMap, cidrList); err != nil {
		fmt.Println("Error loading IPv6 file:", err)
		return
	}

	geoIPList := new(router.GeoIPList)
	for cc, cidr := range cidrList {
		geoIPList.Entry = append(geoIPList.Entry, &router.GeoIP{
			CountryCode: cc,
			Cidr:        cidr,
		})
	}

	geoIPList.Entry = appendExtra(geoIPList.Entry)
	geoIPBytes, err := proto.Marshal(geoIPList)
	if err != nil {
		fmt.Println("Error marshalling geoip list:", err)
	}

	if err := ioutil.WriteFile("geoip.dat", geoIPBytes, 0777); err != nil {
		fmt.Println("Error writing geoip to file:", err)
	}
}

type List struct {
	Name  string
	CIDRList []string
}

var ipKindMap = map[string][]string{}

func DetectPath(path string) (string, error) {
	arrPath := strings.Split(path, string(filepath.ListSeparator))
	for _, content := range arrPath {
		fullPath := filepath.Join(content, "src", "github.com", "yuanmomo", "geoip", "data")
		_, err := os.Stat(fullPath)
		if err == nil || os.IsExist(err) {
			return fullPath, nil
		}
	}
	err := errors.New("No file found in GOPATH")
	return "", err
}

func loadFromFile(){
	dir, err := DetectPath(os.Getenv("GOPATH"))
	if err != nil {
		fmt.Println("Failed to find GOPATH: ", err)
		return
	}

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		list, err := Load(path)
		if err != nil {
			return err
		}
		ipKindMap[list.Name] = list.CIDRList
		return nil
	})
	if err != nil {
		fmt.Println("Failed to load from local file: ", err)
		return
	}
}

func loadFromOnline(){
	cidrMap := list.GetAll()
	for k, v := range cidrMap {
		ipKindMap[k] = v
	}

}

func Load(path string) (*List, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	list := &List{
		Name: strings.ToUpper(filepath.Base(path)),
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		line = removeComment(line)
		if len(line) == 0 {
			continue
		}
		if err != nil {
			return nil, err
		}
		list.CIDRList = append(list.CIDRList, strings.TrimSpace(line))
	}

	return list, nil
}


func removeComment(line string) string {
	idx := strings.Index(line, "#")
	if idx == -1 {
		return line
	}
	return strings.TrimSpace(line[:idx])
}


func appendExtra(geoIPList []*router.GeoIP) []*router.GeoIP{
	// 加载 data 目录
	loadFromFile()
	loadFromOnline()

	for k, v := range ipKindMap {
		cidr := make([]*router.CIDR, 0, 16)
		for _, ip := range v {
			c, err := conf.ParseIP(ip)
			common.Must(err)
			cidr = append(cidr, c)
		}
		geoIP := &router.GeoIP{
			CountryCode: k,
			Cidr:        cidr,
		}
		geoIPList = append(geoIPList, geoIP)
	}
	return  geoIPList
}

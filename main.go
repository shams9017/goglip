package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"strings"

	"github.com/oschwald/geoip2-golang"
)

func displayTopTenResults(countries []string) {
	count := 0

	for c := range countries {
		if count == 10 {
			break
		}
		fmt.Println(countries[c])

		count += 1

	}
}

func extractLineInfo(line string, ipAddress *string, getRequest *string, url *string) {

	parts := strings.Split(line, " ")

	*ipAddress = parts[0]
	*getRequest = parts[3]
	*url = parts[4]

}

func sortCountryNamesDesc(countryFreqMap map[string]int) []string {
	keys := make([]string, 0, len(countryFreqMap))

	for key := range countryFreqMap {
		// if key == "" {
		// 	continue
		// }
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return countryFreqMap[keys[i]] > countryFreqMap[keys[j]]
	})

	//fmt.Println(keys)
	return keys
}

func main() {

	ipAddress := ""
	getRequest := ""
	url := ""
	file, err := os.Open("input2.txt")
	countryNames := []string{}
	cityNames := []string{}
	countryFreqMap := make(map[string]int)

	scanner := bufio.NewScanner(file)

	if err != nil {
		panic(err)
	}
	defer file.Close()

	db, err := geoip2.Open("GeoLite2-City.mmdb")

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Iterate over each line in the file.
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, "images") {
			extractLineInfo(line, &ipAddress, &getRequest, &url)
			ip := net.ParseIP(ipAddress)
			record, err := db.City(ip)
			if err != nil {
				log.Panic(err)
			}
			country := record.Country.Names["en"]
			city := record.City.Names["en"]
			countryNames = append(countryNames, country)
			cityNames = append(cityNames, city)
		}

	}
	for _, c := range countryNames {
		countryFreqMap[c]++
	}
	//fmt.Println(countryFreqMap)
	keys := sortCountryNamesDesc(countryFreqMap)

	fmt.Println("Top 10 Countries where most visitors are: ")
	displayTopTenResults(keys)

}

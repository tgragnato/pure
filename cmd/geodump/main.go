package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/oschwald/maxminddb-golang"
)

func main() {
	ipv4Only := flag.Bool("4", false, "Show only IPv4")
	ipv6Only := flag.Bool("6", false, "Show only IPv6")
	flag.Parse()

	db, err := maxminddb.Open("GeoLite2-Country.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var record struct {
		Continent struct {
			Code string `maxminddb:"code"`
		} `maxminddb:"continent"`
	}

	networks := db.Networks(maxminddb.SkipAliasedNetworks)
	for networks.Next() {
		subnet, err := networks.Network(&record)
		if err != nil {
			log.Fatal(err)
		}

		isIPv4 := subnet.IP.To4() != nil
		if (*ipv4Only && !isIPv4) || (*ipv6Only && isIPv4) {
			continue
		}

		continent := record.Continent.Code
		if continent != "EU" && continent != "NA" && continent != "" {
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("  - %s\n", subnet)
		}
	}
}

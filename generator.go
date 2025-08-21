package main

import (
	"encoding/csv"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
)

func main() {
	writer, err := mmdbwriter.New(
		mmdbwriter.Options{
			DatabaseType:            "GeoLite2-ASN",
			RecordSize:              24,
			IncludeReservedNetworks: true,
			Description: map[string]string{
				"en": "GeoLite2 ASN data for DN42. Learn more at https://github.com/rdp-studio/dn42-geoasn",
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range []string{"GeoLite2-ASN-DN42-Source.csv"} {
		fh, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}

		r := csv.NewReader(fh)

		// first line
		r.Read()

		for {
			row, err := r.Read()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			if len(row) != 3 {
				log.Fatalf("unexpected CSV rows: %v", row)
			}

			_, network, err := net.ParseCIDR(row[0])
			if err != nil {
				log.Fatal(err)
			}

			asn, err := strconv.Atoi(row[1])
			if err != nil {
				log.Fatal(err)
			}

			record := mmdbtype.Map{}

			if asn != 0 {
				record["autonomous_system_number"] = mmdbtype.Uint32(asn)
			}

			if row[2] != "" {
				record["autonomous_system_organization"] = mmdbtype.String(row[2])
			}

			err = writer.Insert(network, record)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	fh, err := os.Create("GeoLite2-ASN-DN42.mmdb")
	if err != nil {
		log.Fatal(err)
	}

	_, err = writer.WriteTo(fh)
	if err != nil {
		log.Fatal(err)
	}
}

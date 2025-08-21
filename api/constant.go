package main

import "time"

const (
	mmdbURL        = "https://github.com/rdp-studio/dn42-geoasn/releases/latest/download/GeoLite2-ASN-DN42.mmdb"
	mmdbMirrorURL  = "https://gh-proxy.com/" + mmdbURL
	localFilePath  = "./GeoLite2-ASN-DN42.mmdb"
	updateInterval = 6 * time.Hour
)

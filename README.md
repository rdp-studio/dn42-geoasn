# DN42-GeoASN

> **语言 / 語言**: [简体中文](docs/README.zh-CN.md) | [繁體中文](docs/README.zh-TW.md)

A tool to generate GeoLite2-compatible ASN databases for DN42 networks. This project extracts route and ASN information from the DN42 registry and converts it into a MaxMind MMDB format that can be used with GeoIP libraries.

## Overview

DN42-GeoASN consists of two main components:
- **Python finder script** (`finder.py`): Parses the DN42 registry to extract route and origin ASN information
- **Go generator** (`generator.go`): Converts the extracted data into a MaxMind MMDB database file

## Features

- ✅ Supports both IPv4 and IPv6 routes
- ✅ Extracts ASN names from the DN42 registry
- ✅ Generates MaxMind MMDB format compatible with existing GeoIP libraries
- ✅ Handles DN42 registry structure automatically
- ✅ Skips routes without proper ASN names

## Prerequisites

- Python 3.x
- Go 1.24.6 or later
- DN42 registry (cloned locally)

## Quick Start (Pre-built Database)

If you just want to use the DN42 ASN database without building it yourself, you can download the latest pre-built MMDB file:

### Download Latest Release

The latest `GeoLite2-ASN-DN42.mmdb` file is automatically built and released at:
**https://github.com/rdp-studio/dn42-geoasn/releases**

```bash
# Download the latest release
wget https://github.com/rdp-studio/dn42-geoasn/releases/latest/download/GeoLite2-ASN-DN42.mmdb

# Or use curl
curl -LO https://github.com/rdp-studio/dn42-geoasn/releases/latest/download/GeoLite2-ASN-DN42.mmdb
```

The database is automatically updated with the latest DN42 registry data and released regularly.

## Building from Source

If you want to build the database yourself or contribute to the project:

### Installation

1. Clone this repository:
   ```bash
   git clone https://github.com/rdp-studio/dn42-geoasn.git
   cd dn42-geoasn
   ```

2. Clone the DN42 registry:
   ```bash
   git clone https://git.dn42.dev/dn42/registry.git
   ```

3. Install Go dependencies:
   ```bash
   go mod download
   ```

## Usage

### Step 1: Extract Route Data

Run the Python finder script to extract route and ASN information from the DN42 registry:

```bash
python finder.py
```

This will:
- Parse all route and route6 objects in the DN42 registry
- Extract origin ASNs for each route
- Look up ASN names from the aut-num objects
- Generate `GeoLite2-ASN-DN42-Source.csv` with the extracted data

### Step 2: Generate MMDB Database

Run the Go generator to convert the CSV data into an MMDB file:

```bash
go run generator.go
```

This will create `GeoLite2-ASN-DN42.mmdb`, a MaxMind-compatible database file.

### Complete Workflow

```bash
# Extract data from DN42 registry
python finder.py

# Generate MMDB file
go run generator.go
```

## Output Files

- `GeoLite2-ASN-DN42-Source.csv`: Intermediate CSV file containing route, ASN, and organization data
- `GeoLite2-ASN-DN42.mmdb`: Final MaxMind MMDB database file

## CSV Format

The intermediate CSV file contains three columns:
1. **Network**: CIDR notation (e.g., `172.20.0.0/14`)
2. **ASN**: Autonomous System Number (without "AS" prefix)
3. **Organization**: ASN name/organization

## MMDB Structure

The generated MMDB file contains records with the following structure:
- `autonomous_system_number`: ASN as uint32
- `autonomous_system_organization`: Organization name as string

## Usage with GeoIP Libraries

The generated MMDB file can be used with standard MaxMind GeoIP libraries:

### Python (maxminddb)
```python
import maxminddb

with maxminddb.open_database('GeoLite2-ASN-DN42.mmdb') as reader:
    result = reader.get('172.20.0.1')
    print(f"ASN: {result['autonomous_system_number']}")
    print(f"Org: {result['autonomous_system_organization']}")
```

### Go (maxminddb-golang)
```go
package main

import (
    "fmt"
    "log"
    "net"
    
    "github.com/oschwald/maxminddb-golang"
)

type ASNResult struct {
    ASN          uint32 `maxminddb:"autonomous_system_number"`
    Organization string `maxminddb:"autonomous_system_organization"`
}

func main() {
    db, err := maxminddb.Open("GeoLite2-ASN-DN42.mmdb")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    ip := net.ParseIP("172.20.0.1")
    var result ASNResult
    err = db.Lookup(ip, &result)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("ASN: %d\n", result.ASN)
    fmt.Printf("Org: %s\n", result.Organization)
}
```

## Configuration

The following constants can be modified in `constant.py`:
- `REGISTRY_PATH`: Path to the DN42 registry directory (default: "registry")
- `SOURCE_OUTPUT`: Output CSV filename (default: "GeoLite2-ASN-DN42-Source.csv")

## Requirements

### Python Dependencies
- Standard library only (os, re, csv)

### Go Dependencies
- `github.com/maxmind/mmdbwriter` v1.0.0
- `github.com/oschwald/maxminddb-golang` v1.12.0 (indirect)
- `go4.org/netipx` v0.0.0-20220812043211-3cc044ffd68d (indirect)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [DN42](https://dn42.dev/) community for maintaining the registry
- [MaxMind](https://www.maxmind.com/) for the MMDB format and libraries
- DN42 registry maintainers and contributors

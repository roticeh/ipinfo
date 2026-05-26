# Roticeh IPInfo

A high-performance, enterprise-grade IP geolocation, ASN, and device fingerprinting microservice written in Go. Built specifically to handle concurrent, high-volume traffic for authentication systems and threat intelligence.

## Background & Legacy Migration

This project is a complete architectural rewrite of the legacy Python-based microservice ([github.com/onurartan/ipinfo](https://github.com/onurartan/ipinfo)). 

The previous Python implementation suffered from critical bottlenecks at scale:
* **File-Lock Errors:** Handling MaxMind `.mmdb` binary databases under heavy concurrent request loads frequently resulted in fatal file-lock collisions.
* **Performance Constraints:** Python's Global Interpreter Lock (GIL) and heavy memory footprint made the service sluggish and unoptimized for real-time authentication flows.

**The Go Advantage:** By migrating to Go, the service now locks the database binaries directly into RAM upon initialization (`geoip2.Open`). This entirely eliminates disk I/O bottlenecks and file-lock issues. Combined with Go's raw concurrency and the Fiber HTTP engine, the service now resolves complex requests in sub-milliseconds with a drastically reduced memory footprint.

## Core Features

* **Subdivision Accuracy:** Resolves location depth down to the exact region and city (e.g., Aydin / Didim).
* **ASN & ISP Intelligence:** Identifies the autonomous system and internet service provider behind the IP.
* **Smart Device Fingerprinting:** Parses User-Agent strings to accurately determine OS, browser, architecture, and device type (PC, Mobile, Bot) via `uaparser`.
* **Dynamic Field Filtering:** Allows endpoints like `/ip/:ipaddress/location` to fetch specific JSON blocks, saving network bandwidth.
* **Zero-Vulnerability Docker Image:** Runs on a Google Distroless minimal image for ultimate security.

---

## Local Development & Testing

To run or test the project on your local machine, we use the minimalist Go build utility, **Craft**, to streamline the development process and handle hot-reloading without bloated Makefiles.

### 1. Install Craft
Ensure you have Go installed on your machine, then install the Craft build engine:
```bash
go install github.com/onurartan/craft@latest

```

*(Make sure your `~/go/bin` path is added to your system environment variables).*

### 2. Clone & Run

Clone the repository and use Craft to start the development server with instant hot-reloading:

```bash
git clone https://github.com/roticeh/ipinfo.git
cd ipinfo

# Starts the server with hot-reload enabled
craft run

```

---

## Production Deployment 🚀

The service is optimized for containerized environments. It uses a multi-stage Docker build, resulting in a lightweight, `distroless` production image that contains zero source code and zero shell vulnerabilities.

### Deploying via Docker Compose

1. Clone the repository to your production server.
2. Ensure you have Docker and Docker Compose installed.
3. Fire up the service:

```bash
docker compose up -d --build

```

### Configuration

The service relies on a `config.yaml` file (or environment variables) to manage server ports, database paths, timeouts, and API rate limits. Refer to the standard `config.yaml` located in the root directory to customize the parameters for your production environment.

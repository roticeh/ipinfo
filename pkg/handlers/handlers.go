package handlers

import (
	"fmt"
	"net"
	"reflect"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/oschwald/geoip2-golang"
	"github.com/ua-parser/uap-go/uaparser"

	cache "github.com/roticeh/ipinfo/pkg/cache"
	utils "github.com/roticeh/ipinfo/pkg/utils"
)

var GeoCityDB *geoip2.Reader
var GeoASNDB *geoip2.Reader
var uaParser *uaparser.Parser

var ipCache *cache.Store[*IPInfoResponse]

type LocationDetail struct {
	Country    string  `json:"country"`
	Continent  string  `json:"continent"`
	Region     string  `json:"region"` // istanbul
	City       string  `json:"city"`   // kadikoy
	PostalCode string  `json:"postal_code"`
	TimeZone   string  `json:"time_zone"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
}

type NetworkDetail struct {
	ASN          uint   `json:"asn"`
	Organization string `json:"organization"`
	ISP          string `json:"isp"`
}

type DeviceDetail struct {
	UserAgentStr string `json:"user_agent_str"`
	BrowserName  string `json:"browser_name"`
	BrowserVer   string `json:"browser_version"`
	OSName       string `json:"os_name"`
	OSVer        string `json:"os_version"`
	OSArch       string `json:"os_architecture"`
	DeviceBrand  string `json:"device_brand"`
	DeviceModel  string `json:"device_model"`
	IsMobile     bool   `json:"is_mobile"`
	IsTablet     bool   `json:"is_tablet"`
	IsPC         bool   `json:"is_pc"`
	IsBot        bool   `json:"is_bot"`
}

type IPInfoResponse struct {
	// Success  bool           `json:"success"`
	// Status   int            `json:"status"`
	ClientIP string         `json:"client_ip"`
	Location LocationDetail `json:"location"`
	Network  NetworkDetail  `json:"network"`
	Device   DeviceDetail   `json:"device"`
}

type BulkIPResponse struct {
	Success bool                   `json:"success"`
	Count   int                    `json:"count"`
	Results map[string]interface{} `json:"results"`
}

func init() {
	uaParser = uaparser.NewFromSaved()

	ipCache = cache.New[*IPInfoResponse](cache.Config{
		TTL:           5 * time.Minute,
		SweepInterval: 10 * time.Minute,
		MaxEntries:    10000,
	})
}

func coreResolveIP(ipStr string, uaStr string) (*IPInfoResponse, error) {

	cacheKey := fmt.Sprintf("%s|%s", ipStr, uaStr)

	if cachedData, found := ipCache.Get(cacheKey); found {
		return cachedData, nil
	}

	parsedIP := net.ParseIP(ipStr)
	if parsedIP == nil {
		return nil, fmt.Errorf("invalid_ip_format")
	}

	if parsedIP.IsLoopback() || parsedIP.IsPrivate() || parsedIP.IsUnspecified() {
		return nil, fmt.Errorf("local_or_private_network_ip")
	}

	// Location Details (Country, Region, City, etc.)
	if GeoCityDB == nil {
		return nil, fmt.Errorf("city_database_not_initialized")
	}
	cityRecord, err := GeoCityDB.City(parsedIP)
	if err != nil {
		return nil, fmt.Errorf("ip_not_found_in_city_database")
	}

	regionName := "Unknown"
	cityName := "Unknown"

	if len(cityRecord.Subdivisions) > 0 {
		if name, ok := cityRecord.Subdivisions[0].Names["en"]; ok {
			regionName = name
		}
	}

	if name, ok := cityRecord.City.Names["en"]; ok && name != "" {
		cityName = name
	} else {
		cityName = regionName
	}

	locDetail := LocationDetail{
		Country:    cityRecord.Country.Names["en"],
		Continent:  cityRecord.Continent.Names["en"],
		Region:     regionName,
		City:       cityName,
		PostalCode: cityRecord.Postal.Code,
		TimeZone:   cityRecord.Location.TimeZone,
		Latitude:   cityRecord.Location.Latitude,
		Longitude:  cityRecord.Location.Longitude,
	}

	// Network Details (ASN, Organization, ISP)
	if GeoASNDB == nil {
		return nil, fmt.Errorf("asn_database_not_initialized")
	}
	asnRecord, err := GeoASNDB.ASN(parsedIP)
	if err != nil {
		return nil, fmt.Errorf("ip_not_found_in_asn_database")
	}

	netDetail := NetworkDetail{
		ASN:          asnRecord.AutonomousSystemNumber,
		Organization: asnRecord.AutonomousSystemOrganization,
		ISP:          asnRecord.AutonomousSystemOrganization,
	}

	// Device Details (User-Agent Analystic)
	devDetail := DeviceDetail{
		UserAgentStr: uaStr,
		BrowserName:  "Unknown",
		BrowserVer:   "Unknown",
		OSName:       "Unknown",
		OSVer:        "Unknown",
		OSArch:       "Unknown",
		DeviceBrand:  "Unknown",
		DeviceModel:  "Unknown",
		IsMobile:     false,
		IsTablet:     false,
		IsPC:         false,
		IsBot:        false,
	}

	if uaStr == "" {
		devDetail.UserAgentStr = "Not Provided"
	}

	if uaStr != "" && uaParser != nil {
		client := uaParser.Parse(uaStr)

		// Browser name and version
		if client.UserAgent.Family != "" && client.UserAgent.Family != "Other" {
			devDetail.BrowserName = client.UserAgent.Family

			browserVer := fmt.Sprintf("%s.%s.%s", client.UserAgent.Major, client.UserAgent.Minor, client.UserAgent.Patch)
			browserVer = strings.TrimSuffix(strings.TrimSuffix(browserVer, "."), ".")
			if browserVer != "" {
				devDetail.BrowserVer = browserVer
			}
		}

		// Operating System Name and Version
		if client.Os.Family != "" && client.Os.Family != "Other" {
			devDetail.OSName = client.Os.Family

			osVer := fmt.Sprintf("%s.%s.%s", client.Os.Major, client.Os.Minor, client.Os.Patch)
			osVer = strings.TrimSuffix(strings.TrimSuffix(osVer, "."), ".")
			if osVer != "" {
				devDetail.OSVer = osVer
			}
		}

		// Device Brand and Model Information
		if client.Device.Brand != "" && client.Device.Brand != "Generic" {
			devDetail.DeviceBrand = client.Device.Brand
		}
		if client.Device.Model != "" {
			devDetail.DeviceModel = client.Device.Model
		}

		// Smart Bot & Device Type Inference
		lowerUA := strings.ToLower(uaStr)
		devDetail.IsBot = client.Device.Family == "Spider" || strings.Contains(lowerUA, "bot") || strings.Contains(lowerUA, "crawl")
		devDetail.IsMobile = strings.Contains(lowerUA, "mobile") || strings.Contains(lowerUA, "android") || strings.Contains(lowerUA, "iphone")
		devDetail.IsTablet = strings.Contains(lowerUA, "ipad") || strings.Contains(lowerUA, "tablet")
		if devDetail.IsTablet {
			devDetail.IsMobile = false
		}

		// if not anything else, assume it's a PC/Desktop
		devDetail.IsPC = !devDetail.IsMobile && !devDetail.IsTablet && !devDetail.IsBot

		// Architecture Detection (x64, arm64, etc.) - This is a heuristic and may not be accurate for all user agents.
		if strings.Contains(lowerUA, "x86_64") || strings.Contains(lowerUA, "win64") || strings.Contains(lowerUA, "x64") {
			devDetail.OSArch = "x64"
		} else if strings.Contains(lowerUA, "arm64") || strings.Contains(lowerUA, "aarch64") {
			devDetail.OSArch = "arm64"
		} else if strings.Contains(lowerUA, "wow64") {
			devDetail.OSArch = "x86"
		}
	}

	response := &IPInfoResponse{
		// Success:  true,
		// Status:   200,
		ClientIP: ipStr,
		Location: locDetail,
		Network:  netDetail,
		Device:   devDetail,
	}

	ipCache.Set(cacheKey, response)
	return response, nil
}

// GetMyIpInfo: request user IP and UA information, analyze and return as JSON response. This is the core API endpoint for client self-analysis. It extracts the client's IP from headers (X-Forwarded-For, X-Real-IP) or falls back to c.IP(). It then calls coreResolveIP to get location, network, and device details.
func GetMyIPInfo(c *fiber.Ctx) error {
	clientIP := c.Get("X-Forwarded-For")
	if clientIP == "" {
		clientIP = c.Get("X-Real-IP")
	}
	if clientIP == "" {
		clientIP = c.IP()
	}

	if strings.Contains(clientIP, ",") {
		clientIP = strings.Split(clientIP, ",")[0]
	}
	clientIP = strings.TrimSpace(clientIP)

	response, err := coreResolveIP(clientIP, c.Get("User-Agent"))
	if err != nil {

		return utils.ResponseError(c, fiber.StatusBadRequest, "request/"+err.Error(), getFriendlyErrorMessage(err.Error()))
	}
	return utils.ResponseSuccess(c, fiber.StatusOK, "", response)
}

// GetSpecificIPInfo: analyze a specific IP address provided as a URL parameter. This endpoint allows users to query information about any IP address, not just their own. It validates the IP format and then calls coreResolveIP to retrieve the details. Error handling ensures that invalid IP formats or private/local IPs are gracefully handled with user-friendly messages.
func GetSpecificIPInfo(c *fiber.Ctx) error {
	targetIP := strings.TrimSpace(c.Params("ipaddress"))
	response, err := coreResolveIP(targetIP, c.Get("User-Agent"))
	if err != nil {
		return utils.ResponseError(c, fiber.StatusBadRequest, "request/"+err.Error(), getFriendlyErrorMessage(err.Error()))
	}
	return utils.ResponseSuccess(c, fiber.StatusOK, "", response)
}

// GetSpecificFieldInfo: extract a specific data block (e.g., /location, /network, /device) dynamically. This endpoint provides granular access to specific sections of the IP information. By passing the desired field as a URL parameter, users can retrieve just the location, network, or device details without the full response. The handler uses reflection to match the requested field against the IPInfoResponse structure and returns only that portion of the data.
func GetSpecificFieldInfo(c *fiber.Ctx) error {
	targetIP := strings.TrimSpace(c.Params("ipaddress"))
	field := strings.ToLower(strings.TrimSpace(c.Params("field")))

	response, err := coreResolveIP(targetIP, c.Get("User-Agent"))
	if err != nil {
		return utils.ResponseError(c, fiber.StatusBadRequest, "request/"+err.Error(), getFriendlyErrorMessage(err.Error()))
	}

	val := reflect.ValueOf(*response)
	typ := reflect.TypeOf(*response)

	for i := 0; i < val.NumField(); i++ {
		jsonTag := typ.Field(i).Tag.Get("json")
		if jsonTag == field {
			return utils.ResponseSuccess(c, fiber.StatusOK, "", fiber.Map{
				"success": true,
				field:     val.Field(i).Interface(),
			})
		}
	}

	return utils.ResponseError(c, fiber.StatusNotFound, "request/field_not_found", "Requested data field not found.")
}

// GetBulkIPInfo: accept a list of IP addresses in the request body and return a map of results for each. This endpoint enables batch processing of multiple IP addresses. It parses the incoming JSON to extract the list of IPs, iterates through them, and calls coreResolveIP for each one. The results are returned in a structured format showing success or failure for each IP.
func GetBulkIPInfo(c *fiber.Ctx) error {
	type BulkIPRequest struct {
		IPs []string `json:"ips"`
	}
	req := new(BulkIPRequest)
	if err := c.BodyParser(req); err != nil {
		return utils.ResponseError(c, fiber.StatusBadRequest, "request/malformed_payload", "Malformed bulk request payload.")
	}

	fmt.Printf("Received bulk IP request for %d IPs\n", len(req.IPs))

	results := make(map[string]interface{})
	for _, ip := range req.IPs {
		ip = strings.TrimSpace(ip)
		response, err := coreResolveIP(ip, "")
		if err != nil {
			results[ip] = fiber.Map{
				"success":    false,
				"error_code": "request/" + err.Error(),
				"message":    getFriendlyErrorMessage(err.Error()),
			}
		} else {
			results[ip] = response
		}
	}

	return utils.ResponseSuccess(c, fiber.StatusOK, "", BulkIPResponse{
		Success: true,
		Count:   len(results),
		Results: results,
	})
}

func getFriendlyErrorMessage(errCode string) string {
	switch errCode {
	case "invalid_ip_format":
		return "The provided string is not a valid IP address."
	case "local_or_private_network_ip":
		return "Loopback or private local network IP addresses cannot be geolocated."
	case "ip_not_found_in_city_database", "ip_not_found_in_asn_database":
		return "The specified IP address does not exist in the public lookup registries."
	default:
		return "An unexpected error occurred during IP resolution."
	}
}

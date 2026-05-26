package tests

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/oschwald/geoip2-golang"
	"github.com/stretchr/testify/assert"

	"github.com/roticeh/ipinfo/pkg/config"
	"github.com/roticeh/ipinfo/pkg/handlers"
	"github.com/roticeh/ipinfo/pkg/routes"
)

func setupTestApp(t *testing.T) *fiber.App {

	config.AppConfig = &config.Config{
		Security: config.SecurityConfig{
			RateLimit: config.RateLimitConfig{
				Max:        150,
				Expiration: 1 * time.Minute,
			},
		},
	}

	config.LoadConfig()

	app := fiber.New()

	var err error
	handlers.GeoCityDB, err = geoip2.Open("../db/GeoLite2-City.mmdb")
	if err != nil {
		t.Fatalf("Failed to load test City database: %v", err)
	}

	handlers.GeoASNDB, err = geoip2.Open("../db/GeoLite2-ASN.mmdb")
	if err != nil {
		t.Fatalf("Failed to load test ASN database: %v", err)
	}

	routes.SetupRoutes(app)
	return app
}

func Test_IPInfo_Routes(t *testing.T) {
	app := setupTestApp(t)

	defer handlers.GeoCityDB.Close()
	defer handlers.GeoASNDB.Close()

	tests := []struct {
		name           string
		method         string
		url            string
		body           string
		reqHeaders     map[string]string
		expectedStatus int
		expectInBody   []string
	}{
		{
			name:           "Valid Public IP Resolution (Aydın IP)",
			method:         "GET",
			url:            "/ip/88.241.68.227",
			body:           "",
			reqHeaders:     map[string]string{"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148.0.0.0 Safari/537.36"},
			expectedStatus: http.StatusOK,
		
			expectInBody: []string{"success\":true", "88.241.68.227", "Aydin", "Türkiye", "Turk Telekom", "Chrome"},
		},
		{
			name:           "Dynamic Field Filtering - Location Block Only",
			method:         "GET",
			url:            "/ip/88.241.68.227/location",
			body:           "",
			reqHeaders:     nil,
			expectedStatus: http.StatusOK,
			expectInBody:   []string{"success\":true", "location", "country", "city", "postal_code"},
		},
		{
			name:           "Sert Kontrol: Reject Localhost IP",
			method:         "GET",
			url:            "/ip/127.0.0.1",
			body:           "",
			reqHeaders:     nil,
			expectedStatus: http.StatusBadRequest,
			expectInBody:   []string{"success\":false", "request/local_or_private_network_ip"},
		},
		{
			name:           "Sert Kontrol: Reject Invalid IP Format",
			method:         "GET",
			url:            "/ip/malformed-ip-string",
			body:           "",
			reqHeaders:     nil,
			expectedStatus: http.StatusBadRequest,
			expectInBody:   []string{"success\":false", "request/invalid_ip_format"},
		},
		{
			name:           "Bulk IP Resolution - Successful Matrix",
			method:         "POST",
			url:            "/ip/bulk",
			body:           `{"ips": ["88.241.68.227", "8.8.8.8"]}`,
			reqHeaders:     map[string]string{"Content-Type": "application/json"},
			expectedStatus: http.StatusOK,
			expectInBody:   []string{"success\":true", "results", "88.241.68.227", "8.8.8.8"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var bodyReader io.Reader
			if tt.body != "" {
				bodyReader = strings.NewReader(tt.body)
			}
			req := httptest.NewRequest(tt.method, tt.url, bodyReader)

			if tt.reqHeaders != nil {
				for k, v := range tt.reqHeaders {
					req.Header.Set(k, v)
				}
			}

			resp, err := app.Test(req, -1)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			respBody, _ := io.ReadAll(resp.Body)
			bodyStr := string(respBody)

			for _, expectedSubstring := range tt.expectInBody {
				assert.Contains(t, bodyStr, expectedSubstring, "Response body should contain: %s", expectedSubstring)
			}
		})
	}
}

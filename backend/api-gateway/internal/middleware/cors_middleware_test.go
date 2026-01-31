package middleware

import (
	"api-gateway/internal/config"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestCors_BasicFunctionality(t *testing.T) {
	cfg := config.Config{
		// Fiber cors middleware expects a comma-separated string for multiple origins
		ClientExternalURL: "https://example.com,https://app.example.com",
	}

	app := fiber.New()
	app.Use(CORS(cfg))

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)

	// CORS working only when Origin header is set
	req.Header.Set("Origin", "https://example.com")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error making request: %v", err)
	}

	assert.Equal(t, fiber.StatusOK, resp.StatusCode, "Expected status OK")

	// Access-Control-Allow-Origin
	allowOrigin := resp.Header.Get("Access-Control-Allow-Origin")
	assert.Equal(t, "https://example.com", allowOrigin, "CORS Allow-Origin header should match request Origin")

	// Access-Control-Allow-Credentials
	allowCredentials := resp.Header.Get("Access-Control-Allow-Credentials")
	assert.Equal(t, "true", allowCredentials, "CORS Allow-Credentials header should be true")
}

func TestCors_PreflightRequest(t *testing.T) {
	cfg := config.Config{
		ClientExternalURL: "https://example.com",
	}

	app := fiber.New()
	app.Use(CORS(cfg))

	app.Post("/api/data", func(c *fiber.Ctx) error {
		return c.SendString("Data received")
	})

	req := httptest.NewRequest(fiber.MethodOptions, "/api/data", nil)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type,Authorization")

	resp, err := app.Test(req)

	assert.NoError(t, err, "Preflight request should not return an error")

	assert.Equal(t, http.StatusNoContent, resp.StatusCode,
		"Preflight request should return status 204 No Content")

	assert.Equal(t, "https://example.com",
		resp.Header.Get("Access-Control-Allow-Origin"))

	assert.Equal(t, "GET,POST,PUT,DELETE,OPTIONS",
		resp.Header.Get("Access-Control-Allow-Methods"),
		"Should allow specified methods")

	assert.Equal(t, "Content-Type,Authorization",
		resp.Header.Get("Access-Control-Allow-Headers"),
		"Should allow specified headers")
}

func TestCORS_UnauthorizedOrigin(t *testing.T) {
	cfg := config.Config{
		ClientExternalURL: "https://example.com",
	}

	app := fiber.New()
	app.Use(CORS(cfg))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Trying with an unauthorized origin
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "https://evil-site.com") // This origin is not allowed

	resp, err := app.Test(req)

	assert.NoError(t, err)

	// CORS headers should not allow this origin
	allowOrigin := resp.Header.Get("Access-Control-Allow-Origin")
	assert.NotEqual(t, "https://evil-site.com", allowOrigin,
		"Unauthorized origin should not be allowed")
}

func TestCORS_MultipleOrigins(t *testing.T) {
	cfg := config.Config{
		ClientExternalURL: "https://app1.example.com,https://app2.example.com,http://localhost:3000",
	}

	app := fiber.New()
	app.Use(CORS(cfg))
	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	testCases := []struct {
		name            string
		origin          string
		shouldBeAllowed bool
	}{
		{
			name:            "First allowed origin",
			origin:          "https://app1.example.com",
			shouldBeAllowed: true,
		},
		{
			name:            "Second allowed origin",
			origin:          "https://app2.example.com",
			shouldBeAllowed: true,
		},
		{
			name:            "Localhost origin",
			origin:          "http://localhost:3000",
			shouldBeAllowed: true,
		},
		{
			name:            "Unauthorized origin",
			origin:          "https://hacker.com",
			shouldBeAllowed: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.Header.Set("Origin", tc.origin)

			resp, err := app.Test(req)

			assert.NoError(t, err)

			allowOrigin := resp.Header.Get("Access-Control-Allow-Origin")

			if tc.shouldBeAllowed {
				// If origin should be allowed, check it is in the header
				msg := fmt.Sprintf("Origin %s should be allowed", tc.origin)
				assert.Equal(t, tc.origin, allowOrigin, msg)
			} else {
				// If origin should not be allowed, check it is not in the header
				msg := fmt.Sprintf("Origin %s should not be allowed", tc.origin)
				assert.NotEqual(t, tc.origin, allowOrigin, msg)
			}
		})
	}
}

func TestCORS_AllHTTPMethods(t *testing.T) {
	// ПОДГОТОВКА
	cfg := config.Config{
		ClientExternalURL: "https://example.com",
	}

	app := fiber.New()
	app.Use(CORS(cfg))

	app.Get("/test", func(c *fiber.Ctx) error { return c.SendString("GET") })
	app.Post("/test", func(c *fiber.Ctx) error { return c.SendString("POST") })
	app.Put("/test", func(c *fiber.Ctx) error { return c.SendString("PUT") })
	app.Delete("/test", func(c *fiber.Ctx) error { return c.SendString("DELETE") })

	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
	}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/test", nil)
			req.Header.Set("Origin", "https://example.com")

			resp, err := app.Test(req)

			assert.NoError(t, err)
			msg := fmt.Sprintf("Method %s should not error", method)
			assert.Equal(t, http.StatusOK, resp.StatusCode, msg)

			assert.Equal(t, "https://example.com",
				resp.Header.Get("Access-Control-Allow-Origin"),
				"CORS header should be set for method "+method)
		})
	}
}

func TestCORS_ResponseBody(t *testing.T) {
	cfg := config.Config{
		ClientExternalURL: "https://example.com",
	}

	app := fiber.New()
	app.Use(CORS(cfg))

	expectedBody := `{"message":"Hello, World!"}`
	app.Get("/json", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello, World!",
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/json", nil)
	req.Header.Set("Origin", "https://example.com")

	resp, err := app.Test(req)

	assert.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "Should read response body without error")

	// Check that the response body is unchanged
	assert.JSONEq(t, expectedBody, string(body),
		"Body response should match expected JSON")

	assert.Equal(t, "https://example.com",
		resp.Header.Get("Access-Control-Allow-Origin"))
}

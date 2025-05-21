package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	"napi/internal/api"
)

func setupTestApp() *fiber.App {
	app := fiber.New()
	api.SetupRoutes(app) // ⚠️ переконайся, що ця функція підключає handler до /api/4.0/tours/getResults
	return app
}

func Test_Handler_GetResults(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest(http.MethodGet, "/api/4.0/tours/getResults?from=1831&to=43&checkIn=2025-04-18&checkTo=2025-04-20&nights=6&nightsTo=8&people=2", nil)
	req.Header.Set("Accept", "application/json")

	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

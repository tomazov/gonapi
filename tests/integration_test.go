package tests

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"napi/internal/api"
)

func TestGetResults(t *testing.T) {
	app := fiber.New()
	api.SetupRoutes(app)

	req := httptest.NewRequest("GET", "/api/4.0/tours/getResults?from=1831&to=43&checkIn=2025-04-18&checkTo=2025-04-20&nights=6&nightsTo=8&people=2", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

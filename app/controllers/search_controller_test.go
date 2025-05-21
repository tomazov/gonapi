package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"

	"napi/app/controllers"
)

func TestGetToursResults(t *testing.T) {
	app := fiber.New()

	// Підключаємо лише цей контролер
	app.Get("/api/4.0/tours/getResults", controllers.GetToursResults)

	req := httptest.NewRequest(http.MethodGet, "/api/4.0/tours/getResults?from=1831&to=43&checkIn=2025-04-18&checkTo=2025-04-20&nights=6&nightsTo=8&people=2", nil)
	req.Header.Set("Accept", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

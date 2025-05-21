package api

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"

	"napi/app/controllers"
	"napi/app/models"
	"napi/app/queries"
	"napi/pkg/config"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/4.0/tours")
	api.Get("/getResults", handleGetResults)
}

func handleGetResults(c *fiber.Ctx) error {
	// 1. Парсимо запит
	var req models.SearchRequest
	if err := c.QueryParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid query params")
	}

	// 2. Генеруємо унікальний пошуковий ключ
	searchID := req.GenerateSearchID()

	log.Printf("🔍 Incoming search: %s | from: %d to: %d", searchID, req.From, req.To)

	// 3. Перевіряємо кеш (чи вже є workProgress)
	workProgress, found := queries.GetWorkProgress(searchID)
	if found {
		if controllers.IsAllDone(workProgress) {
			// Всі ТО обробилися → вертаємо результат
			results := queries.GetFinalResults(searchID, workProgress)
			return c.JSON(fiber.Map{
				"lastResult": true,
				"results":    results,
			})
		}

		// Чекаємо інші ТО
		return c.JSON(fiber.Map{
			"lastResult": false,
			"results":    workProgress,
		})
	}

	// 4. Запускаємо новий пошук
	toOperators := req.ToOperators
	if len(toOperators) == 0 {
		toOperators = queries.LoadDefaultOperators() // SELECT from DB
	}

	controllers.InitSearchTask(searchID, req, toOperators)

	// 5. Відповідь — очікування
	return c.JSON(fiber.Map{
		"lastResult": false,
		"searchId":   searchID,
		"message":    "Search initialized",
	})
}

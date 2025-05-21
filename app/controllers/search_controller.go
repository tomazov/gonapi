package controllers

import (
	"fmt"
	"time"

	"napi/app/models"
	"napi/internal/cache"
	"napi/internal/mq"
)

func InitSearchTask(searchID string, req models.SearchRequest, toOperators []int) {
	progress := make(models.WorkProgress)

	for _, recID := range toOperators {
		// 1. Готуємо task
		task := models.SearchTask{
			SearchID: searchID,
			RecID:    recID,
			Request:  req,
			TOConfig: models.LoadTOConfig(recID), // завантажує fApi, fToken, fUserPass...
			Timestamp: time.Now(),
		}

		// 2. Пуш у RabbitMQ
		_ = mq.PushSearchTask(task)

		// 3. Формуємо workProgress
		progress[fmt.Sprintf("%d", recID)] = models.OperatorStatus{
			Operator: getTOName(recID),
			Code:     100,
			Status:   "run",
			Hotels:   nil,
			Offers:   nil,
			URL:      task.TOConfig.BookingURL,
		}

		// 4. Статус ТО → run
		cache.SetString(fmt.Sprintf("search:%s:%d:status", searchID, recID), "run", cache.TTL())
	}

	// 5. Записуємо загальний прогрес
	cache.SetJSON(fmt.Sprintf("search:%s:workProgress", searchID), progress, cache.TTL())
}

func IsAllDone(progress models.WorkProgress) bool {
	for _, status := range progress {
		if status.Status != "done" {
			return false
		}
	}
	return true
}

func getTOName(recID int) string {
	// Пізніше можна підтягувати з MySQL TO table
	m := map[int]string{
		2700: "TPG2",
		3306: "Aristeya",
		3357: "JoinUp",
		3344: "AllianceEu",
		3419: "LyuboSvitEu",
	}
	if name, ok := m[recID]; ok {
		return name
	}
	return fmt.Sprintf("TO_%d", recID)
}

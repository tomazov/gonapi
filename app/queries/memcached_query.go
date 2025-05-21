package queries

import (
	"fmt"
	"log"

	"napi/app/models"
	"napi/internal/cache"
)

// Отримати workProgress для всього пошуку
func GetWorkProgress(searchID string) (models.WorkProgress, bool) {
	var progress models.WorkProgress
	err := cache.GetJSON(fmt.Sprintf("search:%s:workProgress", searchID), &progress)
	if err != nil {
		log.Printf("🧊 no workProgress found: %v", err)
		return nil, false
	}
	return progress, true
}

// Отримати дані конкретного ТО
func GetTOData(searchID string, recID int) (interface{}, error) {
	var offers interface{}
	err := cache.GetJSON(fmt.Sprintf("search:%s:%d:data", searchID, recID), &offers)
	if err != nil {
		return nil, err
	}
	return offers, nil
}

// Отримати статус ТО
func GetTOStatus(searchID string, recID int) (string, error) {
	return cache.GetString(fmt.Sprintf("search:%s:%d:status", searchID, recID))
}

// Зібрати фінальний результат для клієнта
func GetFinalResults(searchID string, progress models.WorkProgress) map[string]interface{} {
	results := make(map[string]interface{})

	for recIDStr := range progress {
		var offers interface{}
		recID := recIDStr

		// кожен rec_id — string, але в ключі точно такий
		err := cache.GetJSON(fmt.Sprintf("search:%s:%s:data", searchID, recID), &offers)
		if err == nil {
			results[recID] = offers
		}
	}

	return results
}

package vendors

import (
	"napi/app/models"
)

// AdapterFunc — сигнатура будь-якого TO-адаптера
type AdapterFunc func(req models.SearchRequest, cfg models.TOConfig) (interface{}, error)

// Registry — глобальний мап TO-реків до адаптерів
var Registry = make(map[int]AdapterFunc)

// Register — додає адаптер до реєстру
func Register(recID int, fn AdapterFunc) {
	if _, exists := Registry[recID]; exists {
		panic("⚠️ Duplicate TO adapter rec_id: " + string(recID))
	}
	Registry[recID] = fn
}

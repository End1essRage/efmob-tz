package http

import "encoding/json"

// EndDateUpdate представляет обновление для EndDate с тремя состояниями
type EndDateUpdate struct {
	Present bool    // Было ли поле передано в запросе
	Null    bool    // Если true - нужно установить в null
	Value   *string // Значение если не null
}

// UnmarshalJSON для обработки JSON
func (e *EndDateUpdate) UnmarshalJSON(data []byte) error {
	e.Present = true

	// Проверяем на null
	if string(data) == "null" {
		e.Null = true
		e.Value = nil
		return nil
	}

	// Разбираем строку
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	e.Null = false
	e.Value = &s
	return nil
}

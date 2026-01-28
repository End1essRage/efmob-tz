package utils

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
)

func WriteJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload) // Игнорируем ошибку энкодинга, не критично
}

func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request format",
			"code":  "VALIDATION_ERROR",
		})
		return err
	}
	return nil
}

// ParseQuery парсит query-параметры из запроса r в структуру dst.
// dst должен быть указателем на структуру с тегами `schema:"<name>"`.
func ParseQuery(w http.ResponseWriter, r *http.Request, dst interface{}) bool {
	if dst == nil {
		WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "nil dst",
			"code":  "INTERNAL_ERROR",
		})
		return false
	}

	values := r.URL.Query()
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "invalid dst",
			"code":  "INTERNAL_ERROR",
		})
		return false
	}

	v = v.Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		tag := fieldType.Tag.Get("schema")
		if tag == "" {
			continue
		}

		vals, ok := values[tag]
		if !ok || len(vals) == 0 {
			continue
		}

		strVal := vals[0]

		// Разбираем по типу поля
		switch field.Kind() {
		case reflect.String:
			field.SetString(strVal)
		case reflect.Ptr:
			if field.Type().Elem().Kind() == reflect.String {
				field.Set(reflect.ValueOf(&strVal))
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intVal, err := strconv.ParseInt(strVal, 10, 64)
			if err != nil {
				WriteJSON(w, http.StatusBadRequest, map[string]string{
					"error": "invalid param, should be int",
					"code":  "BAD_QUERY",
				})
				return false
			}
			field.SetInt(intVal)
		case reflect.Bool:
			boolVal, err := strconv.ParseBool(strVal)
			if err != nil {
				WriteJSON(w, http.StatusBadRequest, map[string]string{
					"error": "invalid param, should be bool",
					"code":  "BAD_QUERY",
				})
				return false
			}
			field.SetBool(boolVal)
		default:
			// Игнорируем unsupported типы
		}
	}

	return true
}

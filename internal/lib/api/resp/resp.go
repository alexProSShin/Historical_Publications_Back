package resp

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse представляет ответ с ошибкой
type ErrorResponse struct {
	Errors []ErrorDetail `json:"errors"`         // Детали ошибок
	Meta   interface{}   `json:"meta,omitempty"` // Метаданные (опционально)
}

// ErrorDetail представляет описание одной ошибки
type ErrorDetail struct {
	Field   string `json:"field,omitempty"` // Поле, связанное с ошибкой
	Message string `json:"message"`         // Сообщение об ошибке
}

// MetaData - структура для дополнительных данных
type MetaData map[string]interface{}

// WriteJSON записывает данные в формате JSON в `http.ResponseWriter`.
func WriteJSON(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if response == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(response)
}

// WriteError записывает ответ с одной или несколькими ошибками.
func WriteError(w http.ResponseWriter, statusCode int, errors []ErrorDetail, meta MetaData) {
	if len(errors) == 0 {
		errors = SingleError("an unexpected error occurred")
	}

	resp := ErrorResponse{
		Errors: errors,
		Meta:   meta,
	}

	WriteJSON(w, statusCode, resp)
}

// ErrorDetailList упрощает создание списка ошибок.
func ErrorDetailList(field, message string) []ErrorDetail {
	return []ErrorDetail{
		{
			Field:   field,
			Message: message,
		},
	}
}

// SingleError создает единичную ошибку в виде массива.
func SingleError(message string) []ErrorDetail {
	return []ErrorDetail{
		{
			Field:   "",
			Message: message,
		},
	}
}

// ReadJSON считывает тело запроса в структуру.
func ReadJSON[T any](r *http.Request) (T, error) {
	var body T
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	return body, err
}

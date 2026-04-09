package service

import (
	"encoding/json"
	"strings"

	"gorm.io/datatypes"
)

// ValidationError represents a validation error
type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

func NewValidationError(message string) ValidationError {
	return ValidationError{Message: message}
}

// Utility functions
func getCellValue(row []string, index int) string {
	if index >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[index])
}

func buildExtraData(data map[string]string) datatypes.JSON {
	if data == nil || len(data) == 0 {
		return datatypes.JSON([]byte("{}"))
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return datatypes.JSON([]byte("{}"))
	}
	return datatypes.JSON(bytes)
}

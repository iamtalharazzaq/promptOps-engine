package services

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

// JSONValidator provides methods to validate JSON strings against schemas.
type JSONValidator struct{}

// NewJSONValidator creates a new JSONValidator.
func NewJSONValidator() *JSONValidator {
	return &JSONValidator{}
}

// Validate checks if the given jsonStr conforms to the schemaStr.
// It returns nil if valid, or an error describing the validation failure.
func (v *JSONValidator) Validate(jsonStr, schemaStr string) error {
	schemaLoader := gojsonschema.NewStringLoader(schemaStr)
	documentLoader := gojsonschema.NewStringLoader(jsonStr)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("schema validation error: %w", err)
	}

	if !result.Valid() {
		var errors string
		for _, desc := range result.Errors() {
			errors += fmt.Sprintf("- %s\n", desc)
		}
		return fmt.Errorf("invalid JSON:\n%s", errors)
	}

	return nil
}

package validator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xeipuuv/gojsonschema"
)

// Validate checks if the provided JSON data adheres to the schema defined in schemaPath.
func Validate(jsonData []byte, schemaPath string) error {
	// Ensure the schema path is absolute for the loader
	absSchemaPath, err := filepath.Abs(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for schema: %w", err)
	}

	schemaLoader := gojsonschema.NewReferenceLoader("file://" + absSchemaPath)
	documentLoader := gojsonschema.NewBytesLoader(jsonData)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("schema validation error: %w", err)
	}

	if !result.Valid() {
		var errMsgs string
		for _, desc := range result.Errors() {
			errMsgs += fmt.Sprintf("- %s\n", desc)
		}
		return fmt.Errorf("JSON does not validate against schema %q:\n%s", schemaPath, errMsgs)
	}

	return nil
}

// ValidateFile reads a JSON file and validates it against a schema file.
func ValidateFile(dataPath, schemaPath string) ([]byte, error) {
	data, err := os.ReadFile(dataPath)
	if err != nil {
		return nil, fmt.Errorf("reading data file %q: %w", dataPath, err)
	}

	if err := Validate(data, schemaPath); err != nil {
		return nil, err
	}

	return data, nil
}

package validator

import (
	"fmt"
	"os"

	"github.com/tidwall/gjson"
	"github.com/xeipuuv/gojsonschema"
)

// Varidate some json file with schema that its included
// throw error if file not have schema field or mismatch with schema
//
// filePath path to json file
// ignoreNonSchema if true, ignore file that not have schema field
func ValidateJSONSchema(filePath string, ignoreNonSchema bool) error {
	s, err := getSchemaField(filePath)
	if err != nil {
		if ignoreNonSchema {
			return nil
		}
		return err
	}
	result, err := gojsonschema.Validate(
		gojsonschema.NewReferenceLoader(s),
		gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%s", filePath)),
	)
	if err != nil {
		return err
	}

	if !result.Valid() {
		e := fmt.Sprintf("The document %s is not valid. See errors:\n", filePath)
		for _, desc := range result.Errors() {
			e += fmt.Sprintf("- %s\n", desc)
		}
		return fmt.Errorf(e)
	}
	return nil
}

// Get schema field from json file
// if not found return error
//
// return schema field value
func getSchemaField(filePath string) (string, error) {
	// read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("Read text file: %w", err)
	}

	value := gjson.Get(string(data), "$schema")
	if !value.Exists() {
		return "", fmt.Errorf("No $schema field found")
	}

	return value.String(), nil
}

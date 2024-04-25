package validator

import (
	"fmt"
	"os"

	"github.com/tidwall/gjson"
	"github.com/xeipuuv/gojsonschema"
)

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

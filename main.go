package main

import (
	"fmt"
	"log"
	"os"

	"path/filepath"

	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v2"
	"github.com/xeipuuv/gojsonschema"
)

func validateJSONSchema(filePath string, ignoreNonSchema bool) error {
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
		e := fmt.Sprintln("The document is not valid. see errors :")
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

func main() {
	ignoreNonSchema := false
	app := &cli.App{
		Name:  "json-schema-validator",
		Usage: "validate json file by its $schema field",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "quiet-on-non-schema",
				Usage:       "Be quiet error logging files without $schema field",
				Destination: &ignoreNonSchema,
				Aliases:     []string{"q"},
			},
		},
		Action: func(cCtx *cli.Context) error {
			hasError := false
			for i := 0; i < cCtx.NArg(); i++ {
				f := cCtx.Args().Get(i)
				prefix := fmt.Sprintf("%s: ", f)
				path, err := filepath.Abs(f)
				if err != nil {
					log.Println(prefix + err.Error())
					hasError = true
					continue
				}
				if err = validateJSONSchema(path, ignoreNonSchema); err != nil {
					log.Println(prefix + err.Error())
					continue
				}
				log.Println(prefix + "ok")
			}
			if hasError {
				return fmt.Errorf("Some files are not valid")
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

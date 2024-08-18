package main

import (
	"fmt"
	"log"
	"os"

	"path/filepath"

	validator "github.com/Omochice/json-schema-validator"
	"github.com/urfave/cli/v2"
)

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
				if err = validator.ValidateJSONSchema(path, ignoreNonSchema); err != nil {
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

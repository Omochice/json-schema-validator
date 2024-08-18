package main

import (
	"log"
	"os"

	validator "github.com/Omochice/json-schema-validator"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "json-schema-validator",
		Usage: "validate json file by its $schema field",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "quiet-on-non-schema",
				Usage:   "Be quiet error logging files without $schema field",
				Aliases: []string{"q"},
			},
		},
		Action: func(cCtx *cli.Context) error {
			var meg multierror.Group
			ignoreNonSchema := cCtx.Bool("quiet-on-non-schema")
			for i := 0; i < cCtx.NArg(); i++ {
				meg.Go(func() error {
					return validator.ValidateJSONSchema(cCtx.Args().Get(i), ignoreNonSchema)
				})
			}

			merr := meg.Wait()
			if merr != nil {
				return merr
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

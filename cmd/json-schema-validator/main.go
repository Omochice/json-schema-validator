package main

import (
	"log"
	"os"

	validator "github.com/Omochice/json-schema-validator"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
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
			var g errgroup.Group
			ignoreNonSchema := cCtx.Bool("quiet-on-non-schema")
			for i := 0; i < cCtx.NArg(); i++ {
				g.Go(func() error {
					return validator.ValidateJSONSchema(cCtx.Args().Get(i), ignoreNonSchema)
				})
			}
			if err := g.Wait(); err != nil {
				return err
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"path/filepath"

	validator "github.com/Omochice/json-schema-validator"
	"github.com/urfave/cli/v2"
)

func run(filename string, ignoreNonSchema bool, wg *sync.WaitGroup, errCh chan error) {
	defer wg.Done()
	path, err := filepath.Abs(filename)
	if err != nil {
		log.Println(err)
		errCh <- err
		return
	}
	if err := validator.ValidateJSONSchema(path, ignoreNonSchema); err != nil {
		log.Println(err)
		errCh <- err
	}
}

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
			errCh := make(chan error, cCtx.NArg())
			var wg sync.WaitGroup
			for i := 0; i < cCtx.NArg(); i++ {
				wg.Add(1)
				go run(cCtx.Args().Get(i), cCtx.Bool("quiet-on-non-schema"), &wg, errCh)
			}
			wg.Wait()
			go func() {
				for err := range errCh {
					fmt.Fprintln(os.Stderr, err)
				}
			}()

			close(errCh)

			if len(errCh) > 0 {
				return errors.New("error occurred")
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}

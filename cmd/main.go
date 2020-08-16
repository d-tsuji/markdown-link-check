package main

import (
	"log"
	"os"

	markdownlink "github.com/d-tsuji/markdownlink"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "mlc",
		Usage: "Markdown Link Checker",
		Action: func(c *cli.Context) error {
			return markdownlink.Check(markdownlink.NewConfig(c))
		},
		Flags:   []cli.Flag{},
		Version: markdownlink.Version,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

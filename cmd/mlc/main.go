package main

import (
	"errors"
	"log"
	"os"

	"github.com/d-tsuji/mlc"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "mlc",
		Usage: "Markdown Link Checker",
		Action: func(c *cli.Context) error {
			if c.Bool("all") {
				if c.String("user") == "" {
					return errors.New(`required "--user [user]" flag`)
				}
				if c.String("repo") == "" {
					return errors.New(`required "--repo [repo]" flag`)
				}
				if c.String("branch") == "" {
					return errors.New(`required "--branch [branch]" flag`)
				}
				if c.String("token") == "" {
					return errors.New(`required "--token [token]" flag`)
				}
			} else {
				if c.Args().Len() != 1 {
					return errors.New(`URL must be required or "--all" flag forget to scan repository`)
				}
			}
			return mlc.Run(mlc.NewConfig(c))
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "user",
				Aliases: []string{"u"},
				Usage:   "target GitHub repository user",
			},
			&cli.StringFlag{
				Name:    "repo",
				Aliases: []string{"r"},
				Usage:   "target GitHub repository name",
			},
			&cli.StringFlag{
				Name:    "branch",
				Aliases: []string{"b"},
				Usage:   "target GitHub repository branch",
				Value:   "master",
			},
			&cli.StringFlag{
				Name:    "token",
				Aliases: []string{"t"},
				Usage:   "target GitHub access token",
			},
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Usage:   "all scan repository mode (required GitHub access token)",
			},
		},
		Version: mlc.Version,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

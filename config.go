package markdownlink

import "github.com/urfave/cli/v2"

type Config struct {
	args []string
}

func NewConfig(c *cli.Context) Config {
	return Config{
		args: c.Args().Slice(),
	}
}

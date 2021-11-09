package main

import (
	"github.com/tiennam886/manager/employee/internal/app/httpapi"
	"github.com/tiennam886/manager/employee/internal/config"
	"github.com/tiennam886/manager/employee/internal/persistence"
	"time"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "Employee Management"
	app.Version = "0.0.1"
	app.Compiled = time.Now()
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "env",
			Aliases: []string{"e"},
			Value:   "../configs/.env",
			Usage:   "set path to environment file",
		},
		&cli.StringFlag{
			Name:    "env_prefix",
			Aliases: []string{"p"},
			Value:   "EMP",
			Usage:   "set path to environment prefix",
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:   "serve",
			Usage:  "Start the core server",
			Action: Serve,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "addr",
					Aliases: []string{"address"},
					Value:   "localhost:8080",
					Usage:   "specify which address to serve on",
				},
			},
		},
	}
}

func Serve(c *cli.Context) error {
	if err := config.LoadEnvFromFile(c.String("env_prefix"), c.String("env")); err != nil {
		return err
	}

	if err := persistence.LoadEmployeeRepositoryWithMongoDB(); err != nil {
		return err
	}

	return httpapi.Serve(c.Context, c.String("addr"))
}

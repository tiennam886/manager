package main

import (
	"context"
	"fmt"
	"github.com/tiennam886/manager/team/internal/app/httpapi"
	"github.com/tiennam886/manager/team/internal/config"
	"github.com/tiennam886/manager/team/internal/persistence"
	"github.com/urfave/cli/v2"
	"os"
	"time"
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
			Value:   "../../configs/.env",
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
					Value:   "localhost:8081",
					Usage:   "specify which address to serve on",
				},
			},
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	err := app.RunContext(ctx, os.Args)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func Serve(c *cli.Context) error {
	if err := config.LoadEnvFromFile(c.String("env_prefix"), c.String("env")); err != nil {
		return err
	}

	if err := persistence.LoadTeamRepositoryWithMongoDB(); err != nil {
		return err
	}

	return httpapi.Serve(c.Context, c.String("addr"))
}

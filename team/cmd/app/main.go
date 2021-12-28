package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/tiennam886/manager/team/internal/service"

	"github.com/tiennam886/manager/team/internal/app/httpapi"
	"github.com/tiennam886/manager/team/internal/config"
	"github.com/tiennam886/manager/team/internal/persistence"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "Teams Management"
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
			Value:   "TEAM",
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
					Value:   "team:8081",
					Usage:   "specify which address to serve on",
				},
				&cli.StringFlag{
					Name:    "db",
					Aliases: []string{"d"},
					Value:   "mongo",
					Usage:   "set name of database to use",
				},
			},
		},
		{
			Name:        "cli",
			Usage:       "Cli Mode",
			Subcommands: teamCli,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "db",
					Aliases: []string{"d"},
					Value:   "postgres",
					Usage:   "set name of database to use",
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
	if err := Load(c); err != nil {
		return err
	}

	return httpapi.Serve(c.Context, c.String("addr"))
}

func Load(c *cli.Context) error {
	if err := config.LoadEnvFromFile(c.String("env_prefix"), c.String("env")); err != nil {
		return err
	}

	if err := persistence.LoadTeamRepository(c.String("db")); err != nil {
		return err
	}
	return nil
}

var teamCli = []*cli.Command{
	{
		Name:  "add",
		Usage: "add a team. 'add -h' for more help",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "description", Value: "", Aliases: []string{"d"}},
			&cli.StringFlag{Name: "name", Value: "", Aliases: []string{"n"}, Usage: "Input name, ex: 'team A', 'tester B'"},
		},
		Action: addTeam,
	},
	{
		Name:  "delete",
		Usage: "delete a team by ID, 'delete -h' for more help",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "id", Value: ""},
		},
		Action: delTeamCmd,
	},
	{
		Name:  "show-all-members",
		Usage: "show information of a team by id, 'show-all-member -h' for more help",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "id", Value: ""},
		},
		Action: showTeamMembersCmd,
	},
	{
		Name:  "change-name",
		Usage: "change name of a team by their IDs, 'change-name -h' for more help",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "description", Value: "", Aliases: []string{"d"}},
			&cli.StringFlag{Name: "id", Value: ""},
			&cli.StringFlag{Name: "name", Value: ""},
		},
		Action: updateTeam,
	},
}

func addTeam(c *cli.Context) error {
	if err := Load(c); err != nil {
		return err
	}
	payload := service.AddTeamCommand{
		Name:        c.String("name"),
		Description: c.String("description"),
	}

	team, err := service.AddTeam(c.Context, payload)
	if err != nil {
		return err
	}
	fmt.Printf("Add team with name: %s - %s - successfully with id %s\n", team.Name, team.Description, team.UID)
	return nil
}

//func showAllTeam(c *cli.Context) error {
//	err := api.InitHandler(db)
//	if err != nil {
//		return err
//	}
//
//	_, _, err = service.GetAllTeam(page, limit)
//	if err != nil {
//		return err
//	}
//
//	fmt.Println("\nAll of Teams were showed")
//	return nil
//}

func delTeamCmd(c *cli.Context) error {
	if err := Load(c); err != nil {
		return err
	}
	err := service.DeleteTeamByUID(c.Context, service.DeleteTeamByUIDCommand(c.String("id")))
	if err != nil {
		return err
	}

	fmt.Printf("Delete team with ID: %s successfully\n", c.String("id"))
	return nil
}

func showTeamMembersCmd(c *cli.Context) error {
	if err := Load(c); err != nil {
		return err
	}

	team, err := service.FindTeamByUID(c.Context, service.FindTeamByUIDCommand(c.String("id")))
	if err != nil {
		return err
	}

	fmt.Println(team)
	return nil
}

func updateTeam(c *cli.Context) error {
	payload := service.UpdateTeamCommand{
		Name:        c.String("name"),
		Description: c.String("description"),
	}

	err := service.UpdateTeamByUid(c.Context, service.UpdateTeamByUIDCommand(c.String("id")), payload)
	if err != nil {
		return err
	}
	fmt.Printf("Updated team uid=%s", c.String("id"))
	return nil
}

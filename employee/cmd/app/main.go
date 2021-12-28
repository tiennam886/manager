package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/tiennam886/manager/employee/internal/service"

	"github.com/tiennam886/manager/employee/internal/app/httpapi"
	"github.com/tiennam886/manager/employee/internal/config"
	"github.com/tiennam886/manager/employee/internal/persistence"

	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "Employees Management"
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
					Value:   "employee:8082",
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
			Subcommands: employeeCli,
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

var employeeCli = []*cli.Command{
	{
		Name:  "add",
		Usage: "add a employee, 'add -h' for more help",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "name", Value: "", Aliases: []string{"n"}, Usage: "Input name, ex: 'team A', 'tester B'"},
			&cli.StringFlag{Name: "gender", Value: "", Aliases: []string{"g"}, Usage: "Only Input 'male' or 'female'"},
			&cli.StringFlag{Name: "dob", Value: "", Aliases: []string{"d"}, Usage: "Input date of birth as format 'yyyy-mm-dd', ex:'2010-06-30'"},
		},
		Action: addEmployeeCmd,
	},
	{
		Name:  "show",
		Usage: "show an employee, 'show -h' for more help",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "id", Value: ""},
		},
		Action: showEmployee,
	},
	{
		Name:  "delete",
		Usage: "delete a team by ID, 'delete -h' for more help",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "id", Value: ""},
		},
		Action: delEmpCmd,
	},
	{
		Name:  "update",
		Usage: "update information of a team by id, 'update -h' for more help",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "id", Value: ""},
			&cli.StringFlag{Name: "name", Value: "", Aliases: []string{"n"}, Usage: "Input name, ex: 'team A', 'tester B'"},
			&cli.StringFlag{Name: "gender", Value: "", Aliases: []string{"g"}, Usage: "Only Input 'male' or 'female'"},
			&cli.StringFlag{Name: "dob", Value: "", Aliases: []string{"d"}, Usage: "Input date of birth as format 'yyyy-mm-dd', ex:'2010-06-30'"},
		},
		Action: updateEmpCmd,
	},
}

func addEmployeeCmd(c *cli.Context) error {
	if err := Load(c); err != nil {
		return err
	}

	payload := service.AddEmployeeCommand{
		Name:   c.String("name"),
		Gender: c.String("gender"),
		DOB:    c.String("dob"),
	}
	employee, err := service.AddEmployee(c.Context, payload)
	if err != nil {
		return err
	}

	fmt.Printf("Insert employer name: %s, gender: %s, DoB: %s successfully with id %s\n", employee.Name, employee.Gender, employee.DOB, employee.UID)
	return nil
}

func showEmployee(c *cli.Context) error {
	if err := Load(c); err != nil {
		return err
	}

	employee, err := service.FindStaffByUID(c.Context, service.FindEmployeeByUIDCommand(c.String("id")))
	if err != nil {
		return err
	}
	fmt.Printf("Found employer name: %s, gender: %s, DoB: %s successfully with id %s\n", employee.Name, employee.Gender, employee.DOB, employee.UID)

	return nil
}

func delEmpCmd(c *cli.Context) error {
	if err := Load(c); err != nil {
		return err
	}

	err := service.DeleteEmployeeByUID(c.Context, service.DeleteEmployeeByUIDCommand(c.String("id")))
	if err != nil {
		return err
	}
	fmt.Printf("Employee with id: %s was deleted\n", c.String("id"))
	return nil
}

func updateEmpCmd(c *cli.Context) error {
	if err := Load(c); err != nil {
		return err
	}

	payload := service.UpdateEmployeeCommand{
		Name:   c.String("name"),
		Gender: c.String("gender"),
		DOB:    c.String("dob"),
	}

	err := service.UpdateEmployeeByUid(c.Context, service.UpdateEmployeeByUIDCommand(c.String("id")), payload)
	if err != nil {
		return err
	}

	fmt.Printf("Employee with id: %s was updated\n", c.String("id"))
	return nil
}

func Load(c *cli.Context) error {
	if err := config.LoadEnvFromFile(c.String("env_prefix"), c.String("env")); err != nil {
		return err
	}

	if err := persistence.LoadEmployeeRepository(c.String("db")); err != nil {
		return err
	}
	return nil
}

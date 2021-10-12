package manager

import (
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
)

var addTeam = &cobra.Command{
	Use:   "addTeam",
	Short: "Adding a Team with its name",
	Long: `Adding a Team with its name with cli structure:
manager addTeam NAME
For example: manager addTeam "Team A"
`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			fmt.Println("Please insert only NAME(of team)")
			return
		}

		name, err := validationString(args[0])
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		err = teamMongo.AddTeam(name)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		return

	},
}

var showAllTeam = &cobra.Command{
	Use:   "showAllTeam",
	Short: "Show a list of all Teams",
	Long:  `Show a list of all teams with number of total, page and limit`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		page := 1
		limit := 10
		numArgs := len(args)

		if numArgs != 2 && numArgs != 0 {
			fmt.Println("Too many args, it has only 0 or 2 args")
			return
		}

		if numArgs == 2 {
			page, err = strconv.Atoi(args[0])
			if err != nil {
				page = 1
			}

			limit, err = strconv.Atoi(args[1])
			if err != nil {
				limit = 10
			}
		}

		teams, total, err1 := teamMongo.ShowAllTeam(page, limit)
		if err1 != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Printf("\nList of all Employers in page: %v, limt: %v, total: %v\n", page, limit, total)
		fmt.Printf("ID\t\t\t\tNAME\t\n")
		for i := range teams {
			fmt.Printf("%s\t%s\n", teams[i].ID.Hex(), teams[i].Team)
		}
		fmt.Println("\nAll Employers were showed")

		return
	},
}

var delTeamCmd = &cobra.Command{
	Use:   "delTeam",
	Short: "Deleting a Team by ID",
	Long: `Deleting a Team by ID with structure: 
manager delTeam ID 
For example: manager delTeam 6156b66f75697f7a901022f1`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Please insert team ID to delete ")
			return
		}

		id, err := validationString(args[0])
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		err = teamMongo.DeleteTeamById(id)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		return
	},
}

var showAllTeamMember = &cobra.Command{
	Use:   "showAllTeamMember",
	Short: "Showing all member in a team with id",
	Long: `Showing all member in a team with cli structure:
manager showAllTeamMember TEAM_ID`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Please update by ID ")
			return
		}

		id, err := validationString(args[0])
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		team, err := teamMongo.ShowAllTeamMember(id)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Printf("List Employers in: %s with id: %s\n\n", team[0].Team, team[0].ID.Hex())

		employers := team[0].Member

		fmt.Printf("ID\t\t\t\tNAME\t\tGENDER\tDOB\n")
		for i := range employers {
			fmt.Printf("%s\t%s\t%v\t%s\n", employers[i].ID.Hex(), employers[i].Name, employers[i].Gender, employers[i].DoB)
		}
		fmt.Println("\nAll Employers were showed")

		return

	},
}

var addTeamMember = &cobra.Command{
	Use:   "addTeamMember",
	Short: "adding an employer ID to a Team with its ID",
	Long: `adding an employer ID to a Team with its ID as cli structure:
manager addTeamMember TEAM_ID MEMBER_ID
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("Only required 2 args as: manager addTeamMember TEAM_ID MEMBER_ID ")
			return
		}

		tId, err := validationString(args[0])
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		mId, err := validationString(args[1])
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		err = teamMongo.AddTeamMember(tId, mId)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		return
	},
}

var delTeamMember = &cobra.Command{
	Use:   "delTeamMember",
	Short: "deleting an employer ID to a Team with its ID",
	Long: `deleting an employer ID to a Team with its ID as cli structure:
manager delTeamMember TEAM_ID MEMBER_ID
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("Only required 2 args as: manager delTeamMember TEAM_ID MEMBER_ID ")
			return
		}

		tId, err := validationString(args[0])
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		mId, err := validationString(args[1])
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		err = teamMongo.DelTeamMemberById(tId, mId)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		return
	},
}

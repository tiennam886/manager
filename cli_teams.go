package manager

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addTeam = &cobra.Command{
	Use:   "addTeam",
	Short: "Adding a Team with its name",
	Long: `Adding a Team with its name with cli structure:
app addTeam NAME
For example: app addTeam "Team A"
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

		err = dbAddTeam(name)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	},
}

var showAllTeam = &cobra.Command{
	Use:   "showAllTeam",
	Short: "Show a list of all Teams",
	Long:  `Show a list of all teams with number of total, page and limit`,
	Run: func(cmd *cobra.Command, args []string) {
		page, limit, err := validationArgs(args)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		teams, total, err1 := dbGetAllTeams(page, limit)
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
	},
}

var delTeamCmd = &cobra.Command{
	Use:   "delTeam",
	Short: "Deleting a Team by ID",
	Long: `Deleting a Team by ID with structure: 
app delTeam ID 
For example: app delTeam 6156b66f75697f7a901022f1`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Please insert team ID to delete ")
			return
		}

		err = dbDeleteTeamById(args[0])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	},
}

var showAllTeamMember = &cobra.Command{
	Use:   "showAllTeamMember",
	Short: "Showing all member in a team with id",
	Long: `Showing all member in a team with cli structure:
app showAllTeamMember TEAM_ID`,
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

		team, err := dbShowAllMemberInTeam(id)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Printf("List Employers in: %s with id: %s\n\n", team.Team, team.ID.Hex())

		employers := team.Member

		fmt.Printf("ID\t\t\t\tNAME\t\tGENDER\tDOB\n")
		for i := range employers {
			fmt.Printf("%s\t%s\t%v\t%s\n",
				employers[i].ID.Hex(), employers[i].Name, convertNumToGender(employers[i].Gender), employers[i].DoB)
		}
		fmt.Println("\nAll Employers were showed")
	},
}

var addTeamMember = &cobra.Command{
	Use:   "addTeamMember",
	Short: "adding an employer ID to a Team with its ID",
	Long: `adding an employer ID to a Team with its ID as cli structure:
app addTeamMember TEAM_ID MEMBER_ID
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("Only required 2 args as: app addTeamMember TEAM_ID MEMBER_ID ")
			return
		}

		err = dbAddTeamMember(args[0], args[1])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	},
}

var delTeamMember = &cobra.Command{
	Use:   "delTeamMember",
	Short: "deleting an employer ID to a Team with its ID",
	Long: `deleting an employer ID to a Team with its ID as cli structure:
app delTeamMember TEAM_ID MEMBER_ID
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("Only required 2 args as: app delTeamMember TEAM_ID MEMBER_ID ")
			return
		}

		err = dbDelTeamMemberById(args[0], args[1])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	},
}

var changeTeamName = &cobra.Command{
	Use:   "changeTeamName",
	Short: "Change name of a team with ID",
	Long: `Changing name of a team with its ID as cli structure:
app changeTeamName TEAM_ID NEW_NAME`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("Only required 2 args as: app delTeamMember TEAM_ID MEMBER_ID ")
			return
		}

		err = dbUpdateTeam(args[0], args[1])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	},
}

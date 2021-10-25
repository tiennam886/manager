package manager

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addTeam = &cobra.Command{
	Use:   "addTeam",
	Short: "Adding a MySqlTeam with its name",
	Long: `Adding a MySqlTeam with its name with cli structure:
app addTeam NAME
For example: app addTeam "Team A"
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Please insert only NAME(of team)")
			return
		}

		id, err := dbAddTeam(args[0])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Printf("Add team with name: %s successfully with id %s\n", args[0], id)
	},
}

var showAllTeam = &cobra.Command{
	Use:   "showAllTeam",
	Short: "Show a list of all MongoTeam",
	Long:  `Show a list of all teams with number of total, page and limit`,
	Run: func(cmd *cobra.Command, args []string) {
		page, limit, err := validationArgs(args)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		_, _, err = dbGetAllTeam(page, limit)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("\nAll Employers were showed")
	},
}

var delTeamCmd = &cobra.Command{
	Use:   "delTeam",
	Short: "Deleting a MySqlTeam by ID",
	Long: `Deleting a MySqlTeam by ID with structure: 
app delTeam ID 
For example: app delTeam 6156b66f75697f7a901022f1`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Please insert team ID to delete ")
			return
		}

		err = dbDelTeam(args[0])
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Printf("Delete team with ID: %s successfully\n", args[0])
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

		_, err = dbShowMemberInTeam(id)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("\nAll Employers were showed")
	},
}

var addTeamMember = &cobra.Command{
	Use:   "addTeamMember",
	Short: "adding an employer ID to a MySqlTeam with its ID",
	Long: `adding an employer ID to a MySqlTeam with its ID as cli structure:
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

		fmt.Printf("Employee with ID: %s was added to team with ID: %s\n", args[1], args[0])
	},
}

var delTeamMember = &cobra.Command{
	Use:   "delTeamMember",
	Short: "deleting an employer ID to a MySqlTeam with its ID",
	Long: `deleting an employer ID to a MySqlTeam with its ID as cli structure:
app delTeamMember TEAM_ID MEMBER_ID
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("Only required 2 args as: app delTeamMember TEAM_ID MEMBER_ID ")
			return
		}

		err = dbDelTeamMember(args[0], args[1])
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Printf("Employee with ID: %s was deleted in team with ID: %s\n", args[1], args[0])
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

		err = dbUpdateTeamName(args[0], args[1])
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Printf("Team with ID: %s was changed name to: %s", args[0], args[1])
	},
}

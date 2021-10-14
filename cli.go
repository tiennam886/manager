package manager

import (
	"fmt"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	mode        string
	employeeCol *mongo.Collection
	teamCol     *mongo.Collection
	err         error
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "An Application for Employers Management",
	Long: `This Application allows you to execute CRUD task that is connected to the DB 
and a Server for API`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if mode == "server" {
			err := serverMode()
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println("Server mode")
		}

	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize()

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringVar(&mode, "mode", "", "Set Server mode with --mode=server")

	teamCol, err = connectCol(uri, database, teamCollection)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	employeeCol, err = connectCol(uri, database, employerCollection)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// cli handles employers
	rootCmd.AddCommand(addEmpCmd)
	rootCmd.AddCommand(showAllEmp)
	rootCmd.AddCommand(delEmpCmd)
	rootCmd.AddCommand(updateEmpCmd)

	//cli handles teams
	rootCmd.AddCommand(addTeam)
	rootCmd.AddCommand(showAllTeam)
	rootCmd.AddCommand(delTeamCmd)
	rootCmd.AddCommand(showAllTeamMember)
	rootCmd.AddCommand(addTeamMember)
	rootCmd.AddCommand(delTeamMember)
	rootCmd.AddCommand(changeTeamName)
}

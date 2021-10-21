package manager

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	db   string
	mode string
	err  error

	employeeCol *mongo.Collection
	teamCol     *mongo.Collection
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "An Application for Employers Management",
	Long: `This Application allows you to execute CRUD task that is connected to the DB 
and a Server for API`,
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
	rootCmd.Flags().StringVar(&db, "db", "", "Set database to use, default is Mongo, set to MySql by --db=mysql")

	conf, err = loadConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	teamCol, err = connectCol(conf.MongoTeamsCol)
	if err != nil {
		fmt.Println("a", err)
		return
	}

	employeeCol, err = connectCol(conf.MongoEmployeeCol)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	mySqlDB, err = connectMySql()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	cacheClient = initCache()
	addCmd()
}

func addCmd() {
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

func load() error {

	return nil
}

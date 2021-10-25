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

	err := load()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

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
	conf, err = loadConfig()
	if err != nil {
		return err
	}

	teamCol, err = connectCol(conf.MongoTeamsCol)
	if err != nil {
		return err
	}

	employeeCol, err = connectCol(conf.MongoEmployeeCol)
	if err != nil {
		return err
	}

	mySqlDB, err = connectMySql()
	if err != nil {
		return err
	}

	cacheClient = initCache()
	return nil
}

//func addCmd() {
//	// cli handles employers
//	addEmpCmd.Flags().StringVar(&db, "db", "", "Set database to use, default is Mongo, set to MySql by --db=mysql")
//	rootCmd.AddCommand(addEmpCmd)
//
//	showAllEmp.Flags().StringVar(&db, "db", "", "Set database to use, default is Mongo, set to MySql by --db=mysql")
//	rootCmd.AddCommand(showAllEmp)
//
//	delEmpCmd.Flags().StringVar(&db, "db", "", "Set database to use, default is Mongo, set to MySql by --db=mysql")
//	rootCmd.AddCommand(delEmpCmd)
//
//	updateEmpCmd.Flags().StringVar(&db, "db", "", "Set database to use, default is Mongo, set to MySql by --db=mysql")
//	rootCmd.AddCommand(updateEmpCmd)
//
//	//cli handles teams
//	addTeam.Flags().StringVar(&db, "db", "", "Set database to use, default is Mongo, set to MySql by --db=mysql")
//	rootCmd.AddCommand(addTeam)
//
//	showAllTeam.Flags().StringVar(&db, "db", "", "Set database to use, default is Mongo, set to MySql by --db=mysql")
//	rootCmd.AddCommand(showAllTeam)
//
//	delTeamCmd.Flags().StringVar(&db, "db", "", "Set database to use, default is Mongo, set to MySql by --db=mysql")
//	rootCmd.AddCommand(delTeamCmd)
//
//	showAllTeamMember.Flags().StringVar(&db, "db", "", "Set database to use, default is Mongo, set to MySql by --db=mysql")
//	rootCmd.AddCommand(showAllTeamMember)
//
//	addTeamMember.Flags().StringVar(&db, "db", "", "Set database to use, default is Mongo, set to MySql by --db=mysql")
//	rootCmd.AddCommand(addTeamMember)
//
//	delTeamMember.Flags().StringVar(&db, "db", "", "Set database to use, default is Mongo, set to MySql by --db=mysql")
//	rootCmd.AddCommand(delTeamMember)
//
//	changeTeamName.Flags().StringVar(&db, "db", "", "Set database to use, default is Mongo, set to MySql by --db=mysql")
//	rootCmd.AddCommand(changeTeamName)
//}

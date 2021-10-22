package manager

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addEmpCmd = &cobra.Command{
	Use:   "addEmp",
	Short: "Adding an MongoEmployer with his/her name, gender and DoB",
	Long: `Adding an MongoEmployer with his/her name, gender and DoB with structure: 
app addEmp NAME GENDER DOB
with GENDER: 0 for male and 1 for female, DOB in format yyyy-MM-DD
For example: app addEmp "Tran Nam" 0 "2000-01-01"`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 3 {
			fmt.Println("Please insert by NAME GENDER( 0 for male, 1 for female) BOB(yyyy-MM-DD) ")
			return
		}

		err := dbAddEmployer(args[0], args[1], args[2])
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Printf("Insert employer name: %s, gender: %s, DoB: %s successfully\n", args[0], args[1], args[2])
	},
}

var showAllEmp = &cobra.Command{
	Use:   "showAllEmp",
	Short: "Show a list of all employers",
	Long: `Show a list of all employers with number of total, page and limit with CLI structure
app showAllEmp 1 15`,
	Run: func(cmd *cobra.Command, args []string) {
		page, limit, err := validationArgs(args)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		employers, total, err := dbShowAllEmp(page, limit)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(total, employers)

		//fmt.Printf("\nList of all Employers in page: %v, limit: %v, total: %v\n", page, limit, total)
		//fmt.Printf("ID\t\t\t\tNAME\t\tGENDER\tDOB\n")
		//for i := range employers {
		//	fmt.Printf("%s\t%s\t%v\t%s\n",
		//		employers[i].ID.Hex(), employers[i].Name, employers[i].Gender, employers[i].DoB)
		//}
		fmt.Println("\nAll Employers were showed")
	},
}

var delEmpCmd = &cobra.Command{
	Use:   "delEmp",
	Short: "Deleting an MongoEmployer by ID",
	Long: `Deleting an MongoEmployer by ID with structure: 
app delEmp ID 
For example: app delEmp 6156b66f75697f7a901022f1`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Please delete by ID ")
			return
		}

		err = dbDelEmployee(args[0])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Printf("Employee with id: %s was deleted\n", args[0])
	},
}

var updateEmpCmd = &cobra.Command{
	Use:   "updateEmp",
	Short: "Updating an MongoEmployer by ID with his/her new name, gender and DoB",
	Long: `Adding an MongoEmployer with his/her name, gender and DoB with structure: 
app updateEmp ID NAME GENDER DOB
with GENDER: 0 for male and 1 for female, DOB in format yyyy-MM-DD
For example: app updateEmp 6156b66f75697f7a901022f1 "Tran Nam" 0 "2000-01-01"`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 4 {
			fmt.Println("Please update by ID NAME GENDER( 0 for male, 1 for female) BOB(yyyy-MM-DD) ")
			return
		}

		err = dbUpdateEmployee(args[0], args[1], args[2], args[3])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Printf("MongoEmployer %s was updated:\nName: %s\nGender: %s\nDoB: %s\n", args[0], args[1], args[2], args[3])
	},
}

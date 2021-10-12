package manager

import (
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
	"time"
)

var addEmpCmd = &cobra.Command{
	Use:   "addEmp",
	Short: "Adding an Employer with his/her name, gender and DoB",
	Long: `Adding an Employer with his/her name, gender and DoB with structure: 
manager addEmp NAME GENDER DOB
with GENDER: 0 for male and 1 for female, DOB in format yyyy-MM-DD
For example: manager addEmp "Tran Nam" 0 "2000-01-01"`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 3 {
			fmt.Println("Please insert by NAME GENDER( 0 for male, 1 for female) BOB(yyyy-MM-DD) ")
			return
		}

		var name = args[0]
		var date = args[2]

		gender, err := strconv.Atoi(args[1])
		if err != nil || (gender != 1 && gender != 0) {
			fmt.Println("Please insert GENDER as 0 for male, 1 for female ")
			return
		}

		const layoutISO = "2006-01-02"
		_, err = time.Parse(layoutISO, date)
		if err != nil {
			fmt.Println("Date not in format yyyy-MM-DD")
			return
		}

		err = employerMongo.AddEmployer(name, gender, date)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		return
	},
}

var showAllEmp = &cobra.Command{
	Use:   "showAllEmp",
	Short: "Show a list of all employers",
	Long: `Show a list of all employers with number of total, page and limit with CLI structure
manager showAllEmp 1 15`,
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

		employers, total, err1 := employerMongo.ShowAll(page, limit)
		if err1 != nil {
			fmt.Println(err.Error())
			return
		}
		//for employer : range employers{
		//
		//}
		fmt.Printf("\nList of all Employers in page: %v, limt: %v, total: %v\n", page, limit, total)
		fmt.Printf("ID\t\t\t\tNAME\t\tGENDER\tDOB\n")
		for i := range employers {
			fmt.Printf("%s\t%s\t%v\t%s\n", employers[i].ID.Hex(), employers[i].Name, employers[i].Gender, employers[i].DoB)
		}
		fmt.Println("\nAll Employers were showed")

		return
	},
}

var delEmpCmd = &cobra.Command{
	Use:   "delEmp",
	Short: "Deleting an Employer by ID",
	Long: `Deleting an Employer by ID with structure: 
manager delEmp ID 
For example: manager addEmp 6156b66f75697f7a901022f1`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Please update by ID ")
			return
		}

		var id = args[0]

		err := employerMongo.DeleteEmployer(id)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		return
	},
}

var updateEmpCmd = &cobra.Command{
	Use:   "updateEmp",
	Short: "Updating an Employer by ID with his/her new name, gender and DoB",
	Long: `Adding an Employer with his/her name, gender and DoB with structure: 
manager updateEmp ID NAME GENDER DOB
with GENDER: 0 for male and 1 for female, DOB in format yyyy-MM-DD
For example: manager updateEmp 6156b66f75697f7a901022f1 "Tran Nam" 0 "2000-01-01"`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 4 {
			fmt.Println("Please update by ID NAME GENDER( 0 for male, 1 for female) BOB(yyyy-MM-DD) ")
			return
		}

		var id = args[0]
		var name = args[1]
		var date = args[3]

		gender, err := strconv.Atoi(args[2])
		if err != nil || (gender != 1 && gender != 0) {
			fmt.Println("Please update GENDER as 0 for male, 1 for female ")
			return
		}

		const layoutISO = "2006-01-02"
		_, err = time.Parse(layoutISO, date)
		if err != nil {
			fmt.Println("Date not in format yyyy-MM-DD")
			return
		}

		err = employerMongo.UpdateEmployer(id, name, gender, date)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		return
	},
}

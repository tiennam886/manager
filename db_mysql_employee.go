package manager

import "fmt"

type MySqlEmployee struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
	DoB    string `json:"dob"`
}

var employeeTable = conf.MySqlEmployee

func dbMySqlAddEmployee(name string, gender int, dob string) error {
	qr := fmt.Sprintf("INSERT INTO %s(name, gender, dob) VALUES (?, ?, ?);", employeeTable)
	_, err := mySqlDB.Query(qr, name, gender, dob)
	return err
}

func dbMySqlShowAllEmployees(offset int, limit int) (interface{}, int, error) {
	qr := fmt.Sprintf("SELECT * FROM %s LIMIT ? OFFSET ? ;", employeeTable)
	all, err := mySqlDB.Query(qr, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var employees []MySqlEmployee
	for all.Next() {
		var employee MySqlEmployee
		var gender int
		err = all.Scan(&employee.ID, &employee.Name, &gender, &employee.DoB)
		if err != nil {
			return employees, 0, err
		}
		employee.Gender = convertNumToGender(gender)
		employees = append(employees, employee)
	}
	return employees, len(employees), nil
}

func dbMySqlGetEmployeeByID(id int) (interface{}, error) {
	qr := fmt.Sprintf("SELECT * FROM %s WHERE id=?", employeeTable)
	res := mySqlDB.QueryRow(qr, id)

	var gender int
	var employee MySqlEmployee
	err := res.Scan(&employee.ID, &employee.Name, &gender, &employee.DoB)
	if err != nil {
		return employee, err
	}
	employee.Gender = convertNumToGender(gender)

	return employee, nil
}

func dbMySqlUpdateEmployee(id int, name string, gender int, dob string) error {
	_, err = dbMySqlGetEmployeeByID(id)
	if err != nil {
		return err
	}

	qr := fmt.Sprintf("UPDATE %s SET name=?, gender=?, dob=? WHERE id=?", employeeTable)
	_, err = mySqlDB.Query(qr, name, gender, dob, id)
	return err
}

func dbMySqlDelEmployeeByID(id int) error {
	_, err = dbMySqlGetEmployeeByID(id)
	if err != nil {
		return err
	}

	err = dbMySqlDelTeamMemByMemID(id)
	if err != nil {
		return err
	}

	qr := fmt.Sprintf("DELETE FROM %s WHERE id=?", employeeTable)
	_, err = mySqlDB.Query(qr, id)
	return err
}

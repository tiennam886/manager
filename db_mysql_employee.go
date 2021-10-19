package manager

import "fmt"

func dbMySqlAddEmployee(name string, gender int, dob string) error {
	qr := fmt.Sprintf("INSERT INTO %s(name, gender, dob) VALUES (?, ?, ?);", "employers")
	_, err := mySqlDB.Query(qr, name, gender, dob)
	return err
}

func dbMySqlShowAllEmployees(offset int, limit int) (interface{}, error) {
	qr := fmt.Sprintf("SELECT * FROM %s LIMIT ? OFFSET ? ;", "employers")
	all, err := mySqlDB.Query(qr, limit, offset)
	if err != nil {
		return nil, err
	}

	var employees []EmployerPost
	for all.Next() {
		var employee EmployerPost
		var gender int
		err = all.Scan(&employee.ID, &employee.Name, &gender, &employee.DoB)
		if err != nil {
			return nil, err
		}
		employee.Gender = convertNumToGender(gender)
		employees = append(employees, employee)
	}
	return employees, nil
}

func dbMySqlGetEmployeeByID(id int) (interface{}, error) {
	qr := fmt.Sprintf("SELECT * FROM %s WHERE id=?", "employers")
	res := mySqlDB.QueryRow(qr, id)
	var employer EmployerPost
	var gender int
	err := res.Scan(&employer.ID, &employer.Name, &gender, &employer.DoB)
	if err != nil {
		return nil, err
	}
	employer.Gender = convertNumToGender(gender)
	return employer, nil
}

func dbMySqlUpdateEmployee(id int, name string, gender int, dob string) error {
	_, err = dbMySqlGetEmployeeByID(id)
	if err != nil {
		return err
	}

	qr := fmt.Sprintf("UPDATE %s SET name=?, gender=?, dob=? WHERE id=?", "employers")
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

	qr := fmt.Sprintf("DELETE FROM %s WHERE id=?", "employers")
	_, err = mySqlDB.Query(qr, id)
	return err
}

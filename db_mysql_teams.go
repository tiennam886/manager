package manager

import "fmt"

type Team struct {
	ID   int
	Name string
}

func dbMySqlAddTeam(name string) error {
	qr := fmt.Sprintf("INSERT INTO %s(name) VALUES (?);", "teams")
	_, err := mySqlDB.Query(qr, name)
	return err
}

func dbMySqlShowAllTeams(offset int, limit int) (interface{}, error) {
	qr := fmt.Sprintf("SELECT * FROM %s LIMIT ? OFFSET ? ;", "teams")
	all, err := mySqlDB.Query(qr, limit, offset)
	if err != nil {
		return nil, err
	}

	var teams []Team
	for all.Next() {
		var team Team
		err = all.Scan(&team.ID, &team.Name)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}
	return teams, nil
}

func dbMySqlFindTeamByID(id int) (interface{}, error) {
	qr := fmt.Sprintf("SELECT * FROM %s WHERE id=?", "teams")
	res := mySqlDB.QueryRow(qr, id)
	var team Team
	err := res.Scan(&team.ID, &team.Name)
	if err != nil {
		return nil, err
	}
	return team, nil
}

func dbMySqlUpdateTeam(id int, name string) error {
	_, err = dbMySqlFindTeamByID(id)
	if err != nil {
		return err
	}

	qr := fmt.Sprintf("UPDATE %s SET name=? WHERE id=?", "teams")
	_, err = mySqlDB.Query(qr, name, id)
	return err
}

func dbMySqlDelTeamByID(id int) error {
	_, err = dbMySqlFindTeamByID(id)
	if err != nil {
		return err
	}

	err = dbMySqlDelTeamMemByTeamID(id)
	if err != nil {
		return err
	}

	qr := fmt.Sprintf("DELETE FROM %s WHERE id=?", "teams")
	_, err = mySqlDB.Query(qr, id)
	return err
}

func dbMySqlAddTeamMember(teamId int, memId int) error {
	qr := fmt.Sprintf("INSERT INTO %s(teamId, memId) VALUES(?, ?)", "teamMembers")
	_, err := mySqlDB.Query(qr, teamId, memId)
	return err
}

func dbMySqlDelTeamMember(teamId int, memId int) error {
	qr := fmt.Sprintf("DELETE FROM %s WHERE teamId=? AND memId=?", "teamMembers")
	_, err = mySqlDB.Query(qr, teamId, memId)
	return err
}

func dbMySqlDelTeamMemByTeamID(teamId int) error {
	qr := fmt.Sprintf("DELETE FROM %s WHERE teamId=?", "teamMembers")
	_, err = mySqlDB.Query(qr, teamId)
	return err
}

func dbMySqlDelTeamMemByMemID(memId int) error {
	qr := fmt.Sprintf("DELETE FROM %s WHERE memId=?", "teamMembers")
	_, err = mySqlDB.Query(qr, memId)
	return err
}

func dbMySqlShowTeamMember(teamId int) (interface{}, error) {
	qr := fmt.Sprintf("SELECT * FROM %s WHERE teamId=?", "teamMembers")
	all, err := mySqlDB.Query(qr, teamId)
	var employees []interface{}

	for all.Next() {
		var id, teamId, memId int
		err = all.Scan(&id, &teamId, &memId)
		if err != nil {
			return nil, err
		}
		employee, err := dbMySqlGetEmployeeByID(memId)
		if err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	return employees, nil

}

package manager

import "fmt"

type MySqlTeam struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type MySqlTeamMem struct {
	ID      int             `json:"id"`
	Name    string          `json:"name"`
	Members []MySqlEmployee `json:"members"`
}

func dbMySqlAddTeam(name string) (int64, error) {
	qr := fmt.Sprintf("INSERT INTO %s(name) VALUES (?);", conf.MySqlTeams)
	resp, err := mySqlDB.Exec(qr, name)
	if err != nil {
		return -1, err
	}
	id, err := resp.LastInsertId()
	return id, err
}

func dbMySqlShowAllTeams(offset int, limit int) (interface{}, int, error) {
	qr := fmt.Sprintf("SELECT * FROM %s LIMIT ? OFFSET ? ;", conf.MySqlTeams)
	all, err := mySqlDB.Query(qr, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var teams []MySqlTeam
	for all.Next() {
		var team MySqlTeam
		err = all.Scan(&team.ID, &team.Name)
		if err != nil {
			return nil, 0, err
		}
		teams = append(teams, team)
	}
	return teams, len(teams), nil
}

func dbMySqlGetTeam(id int) (interface{}, error) {
	qr := fmt.Sprintf("SELECT * FROM %s WHERE id=?", conf.MySqlTeams)
	res := mySqlDB.QueryRow(qr, id)

	var team MySqlTeam
	err := res.Scan(&team.ID, &team.Name)
	if err != nil {
		return nil, err
	}

	members, err := dbMySqlShowTeamMember(id)
	if err != nil {
		return nil, err
	}

	data := MySqlTeamMem{
		ID:      team.ID,
		Name:    team.Name,
		Members: members.([]MySqlEmployee),
	}
	return data, nil
}

func dbMySqlUpdateTeam(id int, name string) error {
	_, err = dbMySqlGetTeam(id)
	if err != nil {
		return err
	}

	qr := fmt.Sprintf("UPDATE %s SET name=? WHERE id=?", conf.MySqlTeams)
	_, err = mySqlDB.Query(qr, name, id)
	return err
}

func dbMySqlDelTeamByID(id int) error {
	_, err = dbMySqlGetTeam(id)
	if err != nil {
		return err
	}

	err = dbMySqlDelTeamMemByTeamID(id)
	if err != nil {
		return err
	}

	qr := fmt.Sprintf("DELETE FROM %s WHERE id=?", conf.MySqlTeams)
	_, err = mySqlDB.Query(qr, id)
	return err
}

func dbMySqlAddTeamMember(teamId int, memId int) error {
	qr := fmt.Sprintf("INSERT INTO %s(teamId, memId) VALUES(?, ?)", conf.MySqlTeamMem)
	_, err := mySqlDB.Query(qr, teamId, memId)
	return err
}

func dbMySqlDelTeamMember(teamId int, memId int) error {
	qr := fmt.Sprintf("DELETE FROM %s WHERE teamId=? AND memId=?", conf.MySqlTeamMem)
	_, err = mySqlDB.Query(qr, teamId, memId)
	return err
}

func dbMySqlDelTeamMemByTeamID(teamId int) error {
	qr := fmt.Sprintf("DELETE FROM %s WHERE teamId=?", conf.MySqlTeamMem)
	_, err = mySqlDB.Query(qr, teamId)
	return err
}

func dbMySqlDelTeamMemByMemID(memId int) error {
	qr := fmt.Sprintf("DELETE FROM %s WHERE memId=?", conf.MySqlTeamMem)
	_, err = mySqlDB.Query(qr, memId)
	return err
}

func dbMySqlShowTeamMember(teamId int) (interface{}, error) {
	var employees []MySqlEmployee

	qr := fmt.Sprintf("SELECT * FROM %s WHERE teamId=?", conf.MySqlTeamMem)
	all, err := mySqlDB.Query(qr, teamId)
	if err != nil {
		return employees, err
	}

	for all.Next() {
		var id, teamId, memId int
		err = all.Scan(&id, &teamId, &memId)
		if err != nil {
			return employees, err
		}
		employee, err := dbMySqlGetEmployeeByID(memId)
		if err != nil {
			return employees, err
		}
		employees = append(employees, employee.(MySqlEmployee))
	}
	return employees, nil
}

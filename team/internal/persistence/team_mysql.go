package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/tiennam886/manager/team/internal/config"
	"github.com/tiennam886/manager/team/internal/model"
)

type mysqlTeamRepository struct {
	Database    *sql.DB
	TeamTable   string
	MemberTable string
}

func (m mysqlTeamRepository) FindAll(ctx context.Context, offset int, limit int) ([]model.Team, error) {
	qr := fmt.Sprintf("SELECT * FROM %s LIMIT ? OFFSET ? ;", m.TeamTable)
	all, err := m.Database.Query(qr, limit, (offset-1)*limit)
	if err != nil {
		return []model.Team{}, err
	}

	var teams []model.Team
	for all.Next() {
		var team model.Team
		err = all.Scan(&team.UID, &team.Name, &team.Description)
		if err != nil {
			return teams, err
		}
		teams = append(teams, team)
	}

	return teams, nil
}

func (m mysqlTeamRepository) AddAnEmployee(ctx context.Context, employeeId string, teamId string) error {
	stmt := fmt.Sprintf("insert into %s (employee_id, team_id) values(?, ?);", m.TeamTable)
	_, err := m.Database.Exec(stmt, employeeId, teamId)
	return err
}

func (m mysqlTeamRepository) FindByTeamId(ctx context.Context, teamId string) ([]string, error) {
	qr := fmt.Sprintf("SELECT * FROM %s WHERE team_id=? ;", m.MemberTable)
	all, err := m.Database.Query(qr, teamId)
	if err != nil {
		return nil, err
	}

	var employeeList []string
	for all.Next() {
		var tId, employeeId string
		err = all.Scan(&employeeId, &tId)
		if err != nil {
			return employeeList, err
		}
		employeeList = append(employeeList, employeeId)
	}

	return employeeList, nil
}

func (m mysqlTeamRepository) DeleteByTeamId(ctx context.Context, teamId string) error {
	stmt := fmt.Sprintf("delete from %s where team_id=?", m.TeamTable)
	_, err := m.Database.Exec(stmt, teamId)
	return err
}

func (m mysqlTeamRepository) DeleteAnEmployee(ctx context.Context, employeeId string, teamId string) error {
	stmt := fmt.Sprintf("delete from %s where employee_id=? and team_id=? ;", m.TeamTable)
	_, err := m.Database.Exec(stmt, employeeId, teamId)
	return err
}

func newMySqlTeamRepository() (repo TeamRepository, err error) {
	mySqlDB, err := sql.Open("mysql", config.Get().MysqlUrl)
	if err != nil {
		return
	}

	repo = &mysqlTeamRepository{
		Database:    mySqlDB,
		TeamTable:   config.Get().TeamTable,
		MemberTable: config.Get().TeamMemberTable,
	}
	return
}

func (m mysqlTeamRepository) FindByUID(ctx context.Context, uid string) (model.Team, error) {
	var team model.Team

	stmt := fmt.Sprintf("Select * from %s where uid = %s", m.TeamTable, uid)
	row, err := m.Database.Query(stmt)
	if err != nil {
		return team, err
	}
	row.Next()
	err = row.Scan(&team.UID, &team.Name, &team.Description)
	return team, err
}

func (m mysqlTeamRepository) Save(ctx context.Context, team model.Team) error {
	stmt := fmt.Sprintf("insert into %s (uid, name, description) values(?, ?, ?);", m.TeamTable)
	_, err := m.Database.Exec(stmt, team.UID, team.Name, team.Description)
	return err
}

func (m mysqlTeamRepository) Update(ctx context.Context, uid string, team model.Team) error {
	stmt := fmt.Sprintf("update %s set name=?, description=? where uid=?", m.TeamTable)
	_, err := m.Database.Exec(stmt, team.Name, team.Description, uid)
	return err
}

func (m mysqlTeamRepository) Remove(ctx context.Context, uid string) error {
	stmt := fmt.Sprintf("delete from %s where uid=?", m.TeamTable)
	_, err := m.Database.Exec(stmt, uid)
	return err
}

package persistence

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"github.com/tiennam886/manager/employee/internal/config"
	"github.com/tiennam886/manager/employee/internal/model"
)

type mysqlEmployeeRepository struct {
	Database        *sql.DB
	EmployeeTable   string
	TeamMemberTable string
}

func (m mysqlEmployeeRepository) DeleteFromTeam(ctx context.Context, employeeId string, teamId string) error {
	stmt := fmt.Sprintf("delete from %s where employee_id = ? and team_id = ? ;", m.TeamMemberTable)
	_, err := m.Database.Exec(stmt, employeeId, teamId)
	return err
}

func (m mysqlEmployeeRepository) AddToTeam(ctx context.Context, employeeId string, teamId string) error {
	stmt := fmt.Sprintf("insert into %s (employee_id, team_id) values(?, ?);", m.TeamMemberTable)
	_, err := m.Database.Exec(stmt, employeeId, teamId)
	return err
}

func (m mysqlEmployeeRepository) FindByEmployeeId(ctx context.Context, employeeId string) ([]string, error) {
	qr := fmt.Sprintf("SELECT * FROM %s where employee_id = ? ;", m.TeamMemberTable)
	all, err := m.Database.Query(qr, employeeId)
	if err != nil {
		return nil, err
	}

	var teamList []string
	for all.Next() {
		var teamId, employeeId string
		err = all.Scan(&employeeId, &teamId)
		if err != nil {
			return nil, err
		}
		teamList = append(teamList, teamId)
	}
	return teamList, nil
}

func (m mysqlEmployeeRepository) DeleteByEmployeeId(ctx context.Context, employeeId string) error {
	stmt := fmt.Sprintf("delete from %s where employee_id = ? ;", m.TeamMemberTable)
	_, err := m.Database.Exec(stmt, employeeId)
	return err
}

func (m mysqlEmployeeRepository) FindAll(ctx context.Context, offset int, limit int) ([]model.EmployeePost, error) {
	qr := fmt.Sprintf("SELECT * FROM %s LIMIT ? OFFSET ? ;", m.EmployeeTable)
	all, err := m.Database.Query(qr, limit, (offset-1)*limit)
	if err != nil {
		return []model.EmployeePost{}, err
	}

	var employees []model.EmployeePost
	for all.Next() {
		var employee model.Employee

		err = all.Scan(&employee.UID, &employee.Name, &employee.Gender, &employee.DOB)
		if err != nil {
			return employees, err
		}
		employees = append(employees, employee.ToEmployeePost())
	}
	return employees, nil
}

func newMySqlEmployeeRepository() (repo EmployeeRepository, err error) {
	mySqlDB, err := sql.Open("mysql", config.Get().MysqlUrl)
	if err != nil {
		return
	}

	repo = &mysqlEmployeeRepository{
		Database:        mySqlDB,
		EmployeeTable:   config.Get().EmployeeTable,
		TeamMemberTable: config.Get().TeamMemberTable,
	}
	return
}

func (m mysqlEmployeeRepository) FindByUID(ctx context.Context, uid string) (model.EmployeePost, error) {
	var employee model.Employee

	stmt := fmt.Sprintf("Select * from %s where uid = ?", m.EmployeeTable)
	row, err := m.Database.Query(stmt, uid)
	if err != nil {
		return employee.ToEmployeePost(), err
	}
	if row.Next() {
		err = row.Scan(&employee.UID, &employee.Name, &employee.Gender, &employee.DOB)
		if err != nil {
			return employee.ToEmployeePost(), err
		}
	}

	return employee.ToEmployeePost(), err
}

func (m mysqlEmployeeRepository) Save(ctx context.Context, employee model.Employee) error {
	stmt := fmt.Sprintf("insert into %s (uid, name, gender, dob) values(?, ?, ?, ?);", m.EmployeeTable)
	_, err := m.Database.Exec(stmt, employee.UID, employee.Name, employee.Gender, employee.DOB)
	return err
}

func (m mysqlEmployeeRepository) Update(ctx context.Context, uid string, employee model.Employee) error {
	stmt := fmt.Sprintf("update %s set name=?, gender=? , dob=? where uid=?", m.EmployeeTable)
	_, err := m.Database.Exec(stmt, employee.Name, employee.Gender, employee.DOB, uid)
	return err
}

func (m mysqlEmployeeRepository) Remove(ctx context.Context, uid string) error {
	stmt := fmt.Sprintf("delete from %s where uid=?", m.EmployeeTable)
	_, err := m.Database.Exec(stmt, uid)
	return err
}

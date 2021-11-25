package persistence

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/tiennam886/manager/employee/internal/config"
	"github.com/tiennam886/manager/employee/internal/model"
)

type postgresqlEmployeeRepository struct {
	Database        *sql.DB
	EmployeeTable   string
	TeamMemberTable string
}

func newPostgresqlEmployeeRepository() (repo EmployeeRepository, err error) {
	postgresqlDB, err := sql.Open("postgres", config.Get().PostgresqlUrl)
	if err != nil {
		return
	}

	repo = &postgresqlEmployeeRepository{
		Database:        postgresqlDB,
		EmployeeTable:   config.Get().EmployeeTable,
		TeamMemberTable: config.Get().TeamMemberTable,
	}
	return
}

func (m postgresqlEmployeeRepository) FindAll(ctx context.Context, offset int, limit int) ([]model.EmployeePost, error) {
	qr := fmt.Sprintf("SELECT * FROM %s LIMIT $1 OFFSET $2 ;", m.EmployeeTable)
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

func (m postgresqlEmployeeRepository) FindByUID(ctx context.Context, uid string) (model.EmployeePost, error) {
	var employee model.Employee

	stmt := fmt.Sprintf("Select * from %s where uid = '%s';", m.EmployeeTable, uid)
	row, err := m.Database.Query(stmt)
	if err != nil {
		return employee.ToEmployeePost(), err
	}
	row.Next()
	err = row.Scan(&employee.UID, &employee.Name, &employee.Gender, &employee.DOB)
	return employee.ToEmployeePost(), err
}

func (m postgresqlEmployeeRepository) Save(ctx context.Context, employee model.Employee) error {
	stmt := fmt.Sprintf("insert into %s (uid, name, gender, dob) values($1, $2, $3, $4)", config.Get().EmployeeTable)
	_, err := m.Database.Exec(stmt, employee.UID, employee.Name, employee.Gender, employee.DOB)
	return err
}

func (m postgresqlEmployeeRepository) Update(ctx context.Context, uid string, employee model.Employee) error {
	stmt := fmt.Sprintf("update %s set name=$1, gender=$2 , dob=$3 where uid=$4", m.EmployeeTable)
	_, err := m.Database.Exec(stmt, employee.Name, employee.Gender, employee.DOB, uid)
	return err
}

func (m postgresqlEmployeeRepository) Remove(ctx context.Context, uid string) error {
	stmt := fmt.Sprintf("delete from %s where uid=$1", m.EmployeeTable)
	_, err := m.Database.Exec(stmt, uid)
	return err
}

func (m postgresqlEmployeeRepository) AddToTeam(ctx context.Context, employeeId string, teamId string) error {
	panic("implement me")
}

func (m postgresqlEmployeeRepository) DeleteFromTeam(ctx context.Context, employeeId string, teamId string) error {
	panic("implement me")
}

func (m postgresqlEmployeeRepository) FindByEmployeeId(ctx context.Context, employeeId string) ([]string, error) {
	panic("implement me")
}

func (m postgresqlEmployeeRepository) DeleteByEmployeeId(ctx context.Context, employeeId string) error {
	panic("implement me")
}

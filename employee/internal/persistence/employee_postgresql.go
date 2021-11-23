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
	Database      *sql.DB
	EmployeeTable string
}

func newPostgresqlEmployeeRepository() (repo EmployeeRepository, err error) {
	postgresqlDB, err := sql.Open("postgres", config.Get().PostgresqlUrl)
	if err != nil {
		return
	}

	repo = &postgresqlEmployeeRepository{
		Database:      postgresqlDB,
		EmployeeTable: config.Get().EmployeeTable,
	}
	return
}

func (m postgresqlEmployeeRepository) FindByUID(ctx context.Context, uid string) (model.Employee, error) {
	var employee model.Employee

	stmt := fmt.Sprintf("Select * from %s where uid = '%s';", config.Get().EmployeeTable, uid)
	row, err := m.Database.Query(stmt)
	if err != nil {
		return employee, err
	}
	row.Next()
	err = row.Scan(&employee.UID, &employee.Name, &employee.Gender, &employee.DOB)
	return employee, err
}

func (m postgresqlEmployeeRepository) Save(ctx context.Context, employee model.Employee) error {
	stmt := fmt.Sprintf("insert into %s (uid, name, gender, dob) values($1, $2, $3, $4)", config.Get().EmployeeTable)
	_, err := m.Database.Exec(stmt, employee.UID, employee.Name, employee.Gender, employee.DOB)
	return err
}

func (m postgresqlEmployeeRepository) Update(ctx context.Context, uid string, employee model.Employee) error {
	stmt := fmt.Sprintf("update %s set name=$1, gender=$2 , dob=$3 where uid=$4", config.Get().EmployeeTable)
	_, err := m.Database.Exec(stmt, employee.Name, employee.Gender, employee.DOB, uid)
	return err
}

func (m postgresqlEmployeeRepository) Remove(ctx context.Context, uid string) error {
	stmt := fmt.Sprintf("delete from %s where uid=$1", config.Get().EmployeeTable)
	_, err := m.Database.Exec(stmt, uid)
	return err
}

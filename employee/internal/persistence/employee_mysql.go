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
	Database      *sql.DB
	EmployeeTable string
}

func newMySqlEmployeeRepository() (repo EmployeeRepository, err error) {
	mySqlDB, err := sql.Open("mysql", config.Get().MysqlUrl)
	if err != nil {
		return
	}

	repo = &mysqlEmployeeRepository{
		Database:      mySqlDB,
		EmployeeTable: config.Get().EmployeeTable,
	}
	return
}

func (m mysqlEmployeeRepository) FindByUID(ctx context.Context, uid string) (model.Employee, error) {
	var employee model.Employee

	stmt := fmt.Sprintf("Select * from %s where uid = %s", m.EmployeeTable, uid)
	row, err := m.Database.Query(stmt)
	if err != nil {
		return employee, err
	}
	row.Next()
	err = row.Scan(&employee.UID, &employee.Name, &employee.Gender, &employee.DOB)
	return employee, err
}

func (m mysqlEmployeeRepository) Save(ctx context.Context, employee model.Employee) error {
	stmt := fmt.Sprintf("insert into %s (uid, name, gender, dob) values(?, ?, ?, ?);", config.Get().EmployeeTable)
	_, err := m.Database.Exec(stmt, employee.UID, employee.Name, employee.Gender, employee.DOB)
	return err
}

func (m mysqlEmployeeRepository) Update(ctx context.Context, uid string, employee model.Employee) error {
	stmt := fmt.Sprintf("update %s set name=?, gender=? , dob=? where uid=?", config.Get().EmployeeTable)
	_, err := m.Database.Exec(stmt, employee.Name, employee.Gender, employee.DOB, uid)
	return err
}

func (m mysqlEmployeeRepository) Remove(ctx context.Context, uid string) error {
	stmt := fmt.Sprintf("delete from %s where uid=?", config.Get().EmployeeTable)
	_, err := m.Database.Exec(stmt, uid)
	return err
}

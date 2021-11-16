package persistence

import (
	"context"
	"database/sql"
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
	panic("implement me")
}

func (m mysqlEmployeeRepository) Save(ctx context.Context, employee model.Employee) error {
	panic("implement me")
}

func (m mysqlEmployeeRepository) Update(ctx context.Context, uid string, employee model.Employee) error {
	panic("implement me")
}

func (m mysqlEmployeeRepository) Remove(ctx context.Context, uid string) error {
	panic("implement me")
}

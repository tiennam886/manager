package persistence

import (
	"context"
	"fmt"

	"github.com/tiennam886/manager/employee/internal/model"
)

var (
	employees     EmployeeRepository
	employeeDB    EmployeeRepository
	employeeCache EmployeeRepository
)

type EmployeeRepository interface {
	FindByUID(ctx context.Context, uid string) (model.Employee, error)
	Save(ctx context.Context, employee model.Employee) error
	Update(ctx context.Context, uid string, employee model.Employee) error
	Remove(ctx context.Context, uid string) error
}

type EmployeeRepo struct {
	EmployeeDB    EmployeeRepository
	EmployeeCache EmployeeRepository
}

func (e EmployeeRepo) FindByUID(ctx context.Context, uid string) (model.Employee, error) {
	data, err := e.EmployeeCache.FindByUID(ctx, uid)

	if err == nil {
		return data, err
	}

	data, err = e.EmployeeDB.FindByUID(ctx, uid)

	return data, e.EmployeeCache.Save(ctx, data)
}

func (e EmployeeRepo) Save(ctx context.Context, employee model.Employee) error {
	return e.EmployeeDB.Save(ctx, employee)
}

func (e EmployeeRepo) Update(ctx context.Context, uid string, employee model.Employee) error {
	err := e.EmployeeDB.Update(ctx, uid, employee)
	if err != nil {
		return err
	}

	return e.EmployeeCache.Remove(ctx, uid)
}

func (e EmployeeRepo) Remove(ctx context.Context, uid string) error {
	err := e.EmployeeDB.Remove(ctx, uid)
	if err != nil {
		return err
	}

	return e.EmployeeCache.Remove(ctx, uid)
}

func Employees() EmployeeRepository {
	if employees == nil {
		panic("persistence: employees not initiated")
	}

	return employees
}

func LoadEmployeeRepository(db string) (err error) {
	switch db {
	case "mongo":
		err = LoadEmployeeRepositoryWithMongoDB()
	case "postgres":
		err = LoadEmployeeRepositoryWithPostgresql()
	case "mysql":
		err = LoadEmployeeRepositoryWithMysql()
	default:
		err = fmt.Errorf("invalid database, choose mongo, postgres or mysql")
	}

	if err != nil {
		return err
	}

	err = LoadEmployeeRepositoryWithRedis()
	if err != nil {
		return
	}

	employees = &EmployeeRepo{
		EmployeeDB:    employeeDB,
		EmployeeCache: employeeCache,
	}
	return
}

func LoadEmployeeRepositoryWithMongoDB() (err error) {
	employeeDB, err = newMongoEmployeeRepository()
	return
}

func LoadEmployeeRepositoryWithPostgresql() (err error) {
	employeeDB, err = newPostgresqlEmployeeRepository()
	return
}

func LoadEmployeeRepositoryWithMysql() (err error) {
	employeeDB, err = newMySqlEmployeeRepository()
	return
}

func LoadEmployeeRepositoryWithRedis() (err error) {
	employeeCache, err = newRedisEmployeeRepository()
	return
}

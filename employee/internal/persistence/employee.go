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
	FindAll(ctx context.Context, offset int, limit int) ([]model.EmployeePost, error)
	FindByUID(ctx context.Context, uid string) (model.EmployeePost, error)
	Save(ctx context.Context, employee model.Employee) error
	Update(ctx context.Context, uid string, employee model.Employee) error
	Remove(ctx context.Context, uid string) error

	AddToTeam(ctx context.Context, employeeId string, teamId string) error
	FindByEmployeeId(ctx context.Context, employeeId string) ([]string, error)
	DeleteByEmployeeId(ctx context.Context, employeeId string) error
	DeleteFromTeam(ctx context.Context, employeeId string, teamId string) error
}

type EmployeeRepo struct {
	EmployeeDB    EmployeeRepository
	EmployeeCache EmployeeRepository
}

func (e EmployeeRepo) DeleteFromTeam(ctx context.Context, employeeId string, teamId string) error {
	return e.EmployeeDB.DeleteFromTeam(ctx, employeeId, teamId)
}

func (e EmployeeRepo) AddToTeam(ctx context.Context, employeeId string, teamId string) error {
	return e.EmployeeDB.AddToTeam(ctx, employeeId, teamId)
}

func (e EmployeeRepo) FindByEmployeeId(ctx context.Context, employeeId string) ([]string, error) {
	data, err := e.EmployeeDB.FindByEmployeeId(ctx, employeeId)
	return data, err
}

func (e EmployeeRepo) DeleteByEmployeeId(ctx context.Context, employeeId string) error {
	return e.EmployeeDB.DeleteByEmployeeId(ctx, employeeId)
}

func (e EmployeeRepo) FindAll(ctx context.Context, offset int, limit int) ([]model.EmployeePost, error) {
	data, err := e.EmployeeDB.FindAll(ctx, offset, limit)
	return data, err
}

func (e EmployeeRepo) FindByUID(ctx context.Context, uid string) (model.EmployeePost, error) {
	data, err := e.EmployeeCache.FindByUID(ctx, uid)

	if err == nil {
		return data, err
	}

	data, err = e.EmployeeDB.FindByUID(ctx, uid)

	return data, e.EmployeeCache.Save(ctx, data.ToEmployeeDoc())
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

	err = e.EmployeeDB.DeleteByEmployeeId(ctx, uid)
	if err != nil {
		fmt.Println(err.Error())
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

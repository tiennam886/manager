package persistence

import (
	"context"

	"github.com/tiennam886/manager/employee/internal/model"
)

var employees EmployeeRepository

type EmployeeRepository interface {
	FindByUID(ctx context.Context, uid string) (model.Employee, error)
	Save(ctx context.Context, staff model.Employee) error
	Update(ctx context.Context, uid string, staff model.Employee) error
	Remove(ctx context.Context, uid string) error
}

func Employee() EmployeeRepository {
	if employees == nil {
		panic("persistence: employees not initiated")
	}
	return employees
}

func LoadEmployeeRepositoryWithMongoDB() (err error) {
	employees, err = newMongoStaffRepository()
	return err
}

package service

import (
	"context"
	"github.com/tiennam886/manager/employee/internal/model"
	"github.com/tiennam886/manager/employee/internal/persistence"
)

func GetAllEmployee(ctx context.Context, offset int, limit int) (employees []model.EmployeePost, err error) {
	employees, err = persistence.Employees().FindAll(ctx, offset, limit)
	return
}

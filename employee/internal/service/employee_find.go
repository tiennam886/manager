package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/tiennam886/manager/employee/internal/model"
	"github.com/tiennam886/manager/employee/internal/persistence"
)

type FindEmployeeByUIDCommand string

func (c FindEmployeeByUIDCommand) Valid() error {
	_, err := uuid.Parse(string(c))
	return err
}

func FindStaffByUID(ctx context.Context, command FindEmployeeByUIDCommand) (staff model.EmployeePost, err error) {
	if err = command.Valid(); err != nil {
		return
	}

	return persistence.Employees().FindByUID(ctx, string(command))
}

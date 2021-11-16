package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/tiennam886/manager/employee/internal/persistence"
)

type DeleteEmployeeByUIDCommand string

func (c DeleteEmployeeByUIDCommand) Valid() error {
	_, err := uuid.Parse(string(c))
	return err
}

func DeleteEmployeeByUID(ctx context.Context, command DeleteEmployeeByUIDCommand) error {
	if err := command.Valid(); err != nil {
		return err
	}

	return persistence.Employees().Remove(ctx, string(command))
}

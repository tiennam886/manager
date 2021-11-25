package service

import (
	"context"
	"github.com/tiennam886/manager/employee/internal/persistence"
)

func FindEmployeesTeams(ctx context.Context, command FindEmployeeByUIDCommand) ([]string, error) {
	err := command.Valid()
	if err != nil {
		return nil, err
	}

	data, err := persistence.Employees().FindByEmployeeId(ctx, string(command))
	return data, err
}

package service

import (
	"context"

	"github.com/tiennam886/manager/employee/internal/persistence"
)

func DeleteEmployeeToTeam(ctx context.Context, command EmployeeToTeamCommand) error {
	if err := command.Valid(); err != nil {
		return err
	}

	return persistence.Employees().DeleteFromTeam(ctx, command.EmployeeId, command.TeamId)
}

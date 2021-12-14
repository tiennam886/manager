package service

import (
	"context"
	"github.com/tiennam886/manager/team/internal/persistence"
)

func DeleteEmployeeFromTeam(ctx context.Context, command TeamEmployeeCommand) error {
	if err := command.Valid(); err != nil {
		return err
	}
	return persistence.Teams().DeleteAnEmployee(ctx, command.EmployeeId, command.TeamId)
}

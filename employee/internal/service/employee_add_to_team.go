package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/tiennam886/manager/employee/internal/persistence"
)

type EmployeeToTeamCommand struct {
	EmployeeId string `json:"employee_id"`
	TeamId     string `json:"team_id"`
}

func (a EmployeeToTeamCommand) Valid() error {
	_, err := uuid.Parse(a.EmployeeId)
	if err != nil {
		return err
	}

	_, err = uuid.Parse(a.TeamId)
	return err
}

func AddEmployeeToTeam(ctx context.Context, command EmployeeToTeamCommand) error {
	if err := command.Valid(); err != nil {
		return err
	}

	return persistence.Employees().AddToTeam(ctx, command.EmployeeId, command.TeamId)
}

package service

import (
	"context"
	"github.com/tiennam886/manager/team/internal/persistence"
)

func FindTeamsEmployees(ctx context.Context, command FindTeamByUIDCommand) ([]string, error) {
	err := command.Valid()
	if err != nil {
		return nil, err
	}

	data, err := persistence.Teams().FindByTeamId(ctx, string(command))
	return data, err
}

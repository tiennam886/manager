package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/tiennam886/manager/team/internal/model"
	"github.com/tiennam886/manager/team/internal/persistence"
)

type FindTeamByUIDCommand string

func (c FindTeamByUIDCommand) Valid() error {
	_, err := uuid.Parse(string(c))
	return err
}

func FindTeamByUID(ctx context.Context, command FindTeamByUIDCommand) (staff model.Team, err error) {
	if err = command.Valid(); err != nil {
		return
	}

	return persistence.Teams().FindByUID(ctx, string(command))
}

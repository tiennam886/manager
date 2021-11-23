package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/tiennam886/manager/team/internal/persistence"
)

type DeleteTeamByUIDCommand string

func (c DeleteTeamByUIDCommand) Valid() error {
	_, err := uuid.Parse(string(c))
	return err
}

func DeleteTeamByUID(ctx context.Context, command DeleteTeamByUIDCommand) error {
	if err := command.Valid(); err != nil {
		return err
	}

	return persistence.Teams().Remove(ctx, string(command))
}

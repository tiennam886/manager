package service

import (
	"context"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"

	"github.com/tiennam886/manager/team/internal/model"
	"github.com/tiennam886/manager/team/internal/persistence"
)

type AddTeamCommand struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c AddTeamCommand) Valid() error {
	_, err := govalidator.ValidateStruct(c)
	return err
}

func AddTeam(ctx context.Context, command AddTeamCommand) (team model.Team, err error) {
	if err = command.Valid(); err != nil {
		return
	}

	team = model.Team{
		UID:         uuid.NewString(),
		Name:        command.Name,
		Description: command.Description,
	}
	err = persistence.Teams().Save(ctx, team)
	return team, err
}

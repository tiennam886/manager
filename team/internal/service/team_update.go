package service

import (
	"context"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"

	"github.com/tiennam886/manager/team/internal/model"
	"github.com/tiennam886/manager/team/internal/persistence"
)

type UpdateTeamByUIDCommand string

type UpdateTeamCommand struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c UpdateTeamByUIDCommand) Valid() error {
	_, err := uuid.Parse(string(c))
	return err
}

func UpdateTeamByUid(ctx context.Context, command UpdateTeamByUIDCommand, data UpdateTeamCommand) (err error) {
	if err = command.Valid(); err != nil {
		return err
	}

	if _, err = govalidator.ValidateStruct(data); err != nil {
		return err
	}
	newTeam := model.Team{
		UID:         string(command),
		Name:        data.Name,
		Description: data.Description,
	}

	return persistence.Teams().Update(ctx, string(command), newTeam)
}

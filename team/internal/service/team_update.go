package service

import (
	"context"
	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"github.com/tiennam886/manager/team/internal/model"
	"github.com/tiennam886/manager/team/internal/persistence"
)

type UpdateTeamByUIDCommand string

func (c UpdateTeamByUIDCommand) Valid() error {
	_, err := uuid.Parse(string(c))
	return err
}

func UpdateTeamByUid(ctx context.Context, command UpdateTeamByUIDCommand, data model.Team) (err error) {
	if err = command.Valid(); err != nil {
		return err
	}

	if _, err = govalidator.ValidateStruct(data); err != nil {
		return err
	}

	return persistence.Team().Update(ctx, string(command), data)
}

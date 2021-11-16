package service

import (
	"context"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"

	"github.com/tiennam886/manager/team/internal/model"
	"github.com/tiennam886/manager/team/internal/persistence"
)

type AddTeamCommand struct {
	Name string `json:"name"`
}

func (c AddTeamCommand) Valid() error {
	_, err := govalidator.ValidateStruct(c)
	return err
}

func AddTeam(ctx context.Context, command AddTeamCommand) (employee model.Team, err error) {
	if err = command.Valid(); err != nil {
		return
	}

	employee = model.Team{
		UID:  uuid.NewString(),
		Name: command.Name,
	}
	err = persistence.Team().Save(ctx, employee)
	return employee, err
}

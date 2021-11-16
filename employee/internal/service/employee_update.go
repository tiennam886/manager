package service

import (
	"context"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"

	"github.com/tiennam886/manager/employee/internal/model"
	"github.com/tiennam886/manager/employee/internal/persistence"
)

type UpdateEmployeeByUIDCommand string
type UpdateEmployeeCommand struct {
	Name   string    `json:"name"`
	DOB    time.Time `json:"dob"`
	Gender string    `json:"gender"`
}

func (c UpdateEmployeeByUIDCommand) Valid() error {
	_, err := uuid.Parse(string(c))
	return err
}

func UpdateEmployeeByUid(ctx context.Context, command UpdateEmployeeByUIDCommand, data UpdateEmployeeCommand) (err error) {
	if err = command.Valid(); err != nil {
		return err
	}

	if _, err = govalidator.ValidateStruct(data); err != nil {
		return err
	}
	newEmployee := model.Employee{
		UID:    string(command),
		Name:   data.Name,
		DOB:    data.DOB,
		Gender: ToGenderNum(data.Gender),
	}

	return persistence.Employees().Update(ctx, string(command), newEmployee)
}

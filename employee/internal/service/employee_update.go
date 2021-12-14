package service

import (
	"context"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"

	"github.com/tiennam886/manager/employee/internal/model"
	"github.com/tiennam886/manager/employee/internal/persistence"
)

type UpdateEmployeeByUIDCommand string

type UpdateEmployeeCommand struct {
	Name   string `json:"name" valid:"required"`
	DOB    string `json:"dob"`
	Gender string `json:"gender" valid:"required"`
}

func (c UpdateEmployeeByUIDCommand) Valid() error {
	_, err := uuid.Parse(string(c))
	return err
}

func UpdateEmployeeByUid(ctx context.Context, command UpdateEmployeeByUIDCommand, data UpdateEmployeeCommand) (err error) {
	if err = command.Valid(); err != nil {
		return err
	}

	date, err := ValidateDate(data.DOB)
	if err != nil {
		return
	}

	if _, err = govalidator.ValidateStruct(data); err != nil {
		return err
	}
	newEmployee := model.Employee{
		UID:    string(command),
		Name:   data.Name,
		DOB:    date,
		Gender: ToGenderNum(data.Gender),
	}

	return persistence.Employees().Update(ctx, string(command), newEmployee)
}

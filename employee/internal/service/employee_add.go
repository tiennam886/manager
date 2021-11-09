package service

import (
	"context"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"github.com/tiennam886/manager/employee/internal/model"
)

type AddEmployeeCommand struct {
	Name    string    `json:"name"`
	DOB     time.Time `json:"dob"`
	Gender  string    `json:"gender"`
	Address string    `json:"address"`
}

func (c AddEmployeeCommand) Valid() error {
	_, err := govalidator.ValidateStruct(c)
	return err
}

func AddEmployee(ctx context.Context, command AddEmployeeCommand) (employee model.Employee, err error) {
	if err = command.Valid(); err != nil {
		return
	}

	employee = model.Employee{
		UID:    uuid.NewString(),
		Name:   command.Name,
		DOB:    command.DOB,
		Gender: ToGenderNum(command.Gender),
	}

	return
}

func ToGenderNum(gender string) int {
	genMap := map[string]int{
		"male":   0,
		"female": 1,
	}
	if gender != "male" && gender != "female" {
		return 2
	}
	return genMap[gender]
}

package service

import (
	"context"
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"

	"github.com/tiennam886/manager/employee/internal/model"
	"github.com/tiennam886/manager/employee/internal/persistence"
)

type AddEmployeeCommand struct {
	Name   string `json:"name"`
	DOB    string `json:"dob"`
	Gender string `json:"gender"`
}

func (c AddEmployeeCommand) Valid() error {
	_, err := govalidator.ValidateStruct(c)
	return err
}

func AddEmployee(ctx context.Context, command AddEmployeeCommand) (employeePost model.EmployeePost, err error) {
	if err = command.Valid(); err != nil {
		return
	}

	date, err := ValidateDate(command.DOB)
	if err != nil {
		return
	}

	employee := model.Employee{
		UID:    uuid.NewString(),
		Name:   command.Name,
		DOB:    date,
		Gender: ToGenderNum(command.Gender),
	}
	err = persistence.Employees().Save(ctx, employee)
	employeePost = employee.ToEmployeePost()
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

func ValidateDate(date string) (string, error) {
	const layoutISO = "2006-01-02"
	_, err := time.Parse(layoutISO, date)
	if err != nil {
		return "", fmt.Errorf("date not in format yyyy-MM-DD")
	}

	return date, nil
}

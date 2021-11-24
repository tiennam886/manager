package model

type Employee struct {
	UID    string `json:"uid"`
	Name   string `json:"name"`
	DOB    string `json:"dob"`
	Gender int    `json:"gender"`
}

func (e Employee) ToGenderStr() string {
	if e.Gender != 1 && e.Gender != 0 {
		e.Gender = 2
	}
	genInt := []string{"male", "female", "unknown"}
	return genInt[e.Gender]
}

func (e Employee) ToEmployeePost() EmployeePost {
	return EmployeePost{
		UID:    e.UID,
		Name:   e.Name,
		Gender: e.ToGenderStr(),
		DOB:    e.DOB,
	}
}

type EmployeePost struct {
	UID    string `json:"uid"`
	Name   string `json:"name"`
	DOB    string `json:"dob"`
	Gender string `json:"gender"`
}

func (e EmployeePost) ToEmployeeDoc() Employee {
	genMap := map[string]int{
		"male":    0,
		"female":  1,
		"unknown": 2,
	}
	return Employee{
		UID:    e.UID,
		Name:   e.Name,
		Gender: genMap[e.Gender],
		DOB:    e.DOB,
	}
}

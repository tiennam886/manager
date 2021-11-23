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
	genInt := []string{"male", "female", "invalid"}
	return genInt[e.Gender]
}

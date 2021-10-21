package manager

import (
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func validationDate(date string) (string, error) {
	const layoutISO = "2006-01-02"
	_, err := time.Parse(layoutISO, date)
	if err != nil {
		fmt.Println("Date not in format yyyy-MM-DD")
		return "", err
	}

	return date, nil
}

func validationString(s string) (string, error) {
	if len(s) > 30 || len(s) < 1 {
		return "", fmt.Errorf("length of %s is more than 30 digits or null", s)
	}
	return s, nil
}

func validationGender(gender string) (int, error) {
	genMap := map[string]int{
		"male":   0,
		"female": 1,
	}
	if gender != "male" && gender != "female" {
		return 2, fmt.Errorf("not %s, only male or female", gender)
	}
	return genMap[gender], nil
}

func validationAddEmployer(name string, gender string, date string) (string, int, string, error) {
	n, err := validationString(name)
	if err != nil {
		return "", -1, "", err
	}

	g, err := validationGender(gender)
	if err != nil {
		return "", -1, "", err
	}

	d, err := validationDate(date)
	if err != nil {
		return "", 0, "", err
	}

	return n, g, d, nil

}

func validationObjectID(teamId string, memId string) (primitive.ObjectID, primitive.ObjectID, error) {
	var (
		objId, memberId primitive.ObjectID
		err             error
	)

	teamId, err = validationString(teamId)
	if err != nil {
		return objId, memberId, err
	}

	memId, err = validationString(memId)
	if err != nil {
		return objId, memberId, err
	}

	objId, err = primitive.ObjectIDFromHex(teamId)
	if err != nil {
		return objId, memberId, fmt.Errorf("TEAM_ID %s was invalid\n", teamId)
	}
	memberId, err = primitive.ObjectIDFromHex(memId)
	if err != nil {
		return objId, memberId, fmt.Errorf("Employer_ID %s was invalid", memId)
	}
	return objId, memberId, nil
}

func validationID(teamId string, memId string) (int, int, error) {
	var tId, mId int
	tId, err = strconv.Atoi(teamId)
	if err != nil {
		return 0, 0, err
	}
	mId, err = strconv.Atoi(memId)
	if err != nil {
		return 0, 0, err
	}
	return tId, mId, nil
}

func validationArgs(args []string) (int, int, error) {
	var page, limit int
	numArgs := len(args)

	if numArgs != 2 && numArgs != 0 {
		fmt.Println("Too many args, it has only 0 or 2 args")
		return 1, 10, fmt.Errorf("number of args is %v, it has only 0 or 2 args", numArgs)
	}

	if numArgs == 2 {
		page, err = strconv.Atoi(args[0])
		if err != nil {
			page = 1
		}
		limit, err = strconv.Atoi(args[1])
		if err != nil {
			limit = 10
		}
	}
	return page, limit, nil
}

func convertNumToGender(a int) string {
	genInt := []string{"male", "female", "invalid"}
	return genInt[a]
}

package manager

import (
	"fmt"
	"strconv"
	"time"
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
	g, err := strconv.Atoi(gender)
	if err != nil || (g != 1 && g != 0) {
		return -1, fmt.Errorf("Please insert GENDER as 0 for male, 1 for female ")
	}
	return g, nil
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

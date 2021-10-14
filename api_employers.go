package manager

import (
	"fmt"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

func apiGetEmployers(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10
	}

	employers, total, err := dbShowAllEmployee(page, limit)
	if err != nil {
		responseAllNotFound(c, err)
		return
	}

	last := math.Ceil(float64(total / int64(limit)))
	if last < 1 && total > 0 {
		last = 1
	}

	responseAllEmployeeOK(c, employers, total, page, last, limit)
	return
}

func apiPostEmployer(c *gin.Context) {
	var employer *EmployerPost

	if err := c.BindJSON(&employer); err != nil {
		responseBadRequest(c, "", err)
		return
	}

	name, gender, dob, err := validationAddEmployer(employer.Name, employer.Gender, employer.DoB)
	if err != nil {
		responseBadRequest(c, "", err)
		return
	}

	//insert the newly created object into mongodb
	err = dbAddEmployer(name, gender, dob)
	if err != nil {
		responseInternalServer(c, "", err)
		return
	}

	msg := fmt.Sprintf("Insert employer name: %s, gender: %s, DoB: %s to DB successfully",
		name, convertNumToGender(gender), dob)

	responseEmployerOK(c, employer, msg)
}

func apiDelEmployerByID(c *gin.Context) {
	id, err := validationString(c.Param("id"))
	if err != nil {
		responseBadRequest(c, c.Param("id"), err)
		return
	}

	err = dbDeleteEmployer(id)
	if err != nil {
		responseInternalServer(c, id, err)
		return
	}
	responseOK(c, id, fmt.Sprintf("Employee with ID %s was deleted\n", id))
}

func apiUpdateEmployerByID(c *gin.Context) {
	var employer *EmployerPost

	id, err := validationString(c.Param("id"))
	if err != nil {
		responseBadRequest(c, c.Param("id"), err)
		return
	}

	if err := c.BindJSON(&employer); err != nil {
		responseBadRequest(c, id, err)
		return
	}

	name, gender, dob, err := validationAddEmployer(employer.Name, employer.Gender, employer.DoB)
	if err != nil {
		responseInternalServer(c, id, err)
		return
	}

	err = dbUpdateEmployer(id, name, gender, dob)
	if err != nil {
		responseInternalServer(c, id, err)
		return
	}

	msg := fmt.Sprintf("Employer %s was updated:\nName: %s\nGender: %s\nDoB: %s\n",
		id, name, convertNumToGender(gender), dob)
	responseOK(c, id, msg)
}

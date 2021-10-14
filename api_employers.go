package manager

import (
	"fmt"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetEmployers(c *gin.Context) {
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

func PostEmployer(c *gin.Context) {
	var employer *EmployerPost
	var err error

	if err := c.BindJSON(&employer); err != nil {
		responseBadRequest(c, employer.ID.Hex(), err)
		return
	}

	name, gender, dob, err := validationAddEmployer(employer.Name, employer.Gender, employer.DoB)
	if err != nil {
		responseBadRequest(c, employer.ID.Hex(), err)
		return
	}

	//insert the newly created object into mongodb
	err = dbAddEmployer(name, gender, dob)
	if err != nil {
		responseInternalServer(c, employer.ID.Hex(), err)
		return
	}
	responseEmployerOK(c, employer, "Add Successfully")

	return

}

func DelEmployerByID(c *gin.Context) {
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

func UpdateEmployerByID(c *gin.Context) {
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
	if err := c.BindJSON(&employer); err != nil {
		responseBadRequest(c, employer.ID.Hex(), err)
		return
	}

	err = dbUpdateEmployer(id, name, gender, dob)
	if err != nil {
		responseInternalServer(c, employer.ID.Hex(), err)
		return
	}
	responseEmployerOK(c, employer, "Update successfully\n")
}

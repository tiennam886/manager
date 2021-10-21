package manager

import (
	"fmt"
	"math"
	"strconv"

	json2 "encoding/json"

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

	employers, total, err := dbShowAllEmp(page, limit)
	if err != nil {
		responseAllNotFound(c, err)
		return
	}

	last := math.Ceil(float64(int64(total) / int64(limit)))
	if last < 1 && total > 0 {
		last = 1
	}

	responseAllEmployeeOK(c, employers, int64(total), page, last, limit)
}

func apiGetEmployee(c *gin.Context) {
	id, err := validationString(c.Param("id"))
	if err != nil {
		responseError(c, id, err)
		return
	}
	var employerPost interface{}
	data, _ := getCache(id)
	err = json2.Unmarshal([]byte(data), &employerPost)
	if data != "" && err == nil {
		responseOK(c, employerPost, "Get MongoEmployer successfully")
		return
	}

	employerPost, err = dbGetEmployee(id)
	if err != nil {
		responseError(c, id, err)
		return
	}

	setCache(id, employerPost)
	responseOK(c, employerPost, "Get MongoEmployer successfully")
}

func apiPostEmployer(c *gin.Context) {
	var employer *MongoEmployerPost

	if err := c.BindJSON(&employer); err != nil {
		responseBadRequest(c, "", err)
		return
	}

	err = dbAddEmployer(employer.Name, employer.Gender, employer.DoB)
	if err != nil {
		responseInternalServer(c, "", err)
		return
	}

	msg := fmt.Sprintf("Insert employer name: %s, gender: %s, DoB: %s to DB successfully", employer.Name, employer.Gender, employer.DoB)
	responseOK(c, employer, msg)
}

func apiDelEmployerByID(c *gin.Context) {
	id, err := validationString(c.Param("id"))
	if err != nil {
		responseBadRequest(c, c.Param("id"), err)
		return
	}

	err = dbDelEmployee(id)
	if err != nil {
		responseInternalServer(c, id, err)
		return
	}

	delCache(id)
	responseOK(c, id, fmt.Sprintf("Employee with ID %s was deleted\n", id))
}

func apiUpdateEmployerByID(c *gin.Context) {
	id, err := validationString(c.Param("id"))
	if err != nil {
		responseBadRequest(c, c.Param("id"), err)
		return
	}

	var employerPost MongoEmployerPost
	if err := c.BindJSON(&employerPost); err != nil {
		responseBadRequest(c, id, err)
		return
	}

	err = dbUpdateEmployee(id, employerPost.Name, employerPost.Gender, employerPost.DoB)
	if err != nil {
		responseError(c, id, err)
		return
	}
	delCache(id)
	msg := fmt.Sprintf("MongoEmployer %s was updated:\nName: %s\nGender: %s\nDoB: %s\n",
		id, employerPost.Name, employerPost.Gender, employerPost.DoB)
	responseOK(c, id, msg)
}

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
		responseError(c, employers, err, 404)
		return
	}

	last := math.Ceil(float64(int64(total) / int64(limit)))
	if last < 1 && total > 0 {
		last = 1
	}

	responseAllDataOK(c, employers, int64(total), page, last, limit)
}

func apiGetEmployee(c *gin.Context) {
	id, err := validationString(c.Param("id"))
	if err != nil {
		responseError(c, id, err, 400)
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
		responseError(c, id, err, 404)
		return
	}

	setCache(id, employerPost)

	responseOK(c, employerPost, "Get MongoEmployer successfully")
}

func apiPostEmployer(c *gin.Context) {
	var employer *MongoEmployerPost

	if err := c.BindJSON(&employer); err != nil {
		responseError(c, nil, err, 400)
		return
	}

	err = dbAddEmployer(employer.Name, employer.Gender, employer.DoB)
	if err != nil {
		responseError(c, nil, err, 404)
		return
	}

	msg := fmt.Sprintf("Insert employer name: %s, gender: %s, DoB: %s to DB successfully", employer.Name, employer.Gender, employer.DoB)
	responseOK(c, employer, msg)
}

func apiDelEmployerByID(c *gin.Context) {
	id, err := validationString(c.Param("id"))
	if err != nil {
		responseError(c, nil, err, 400)
		return
	}

	err = dbDelEmployee(id)
	if err != nil {
		responseError(c, nil, err, 404)
		return
	}

	delCache(id)

	responseOK(c, id, fmt.Sprintf("Employee with ID %s was deleted\n", id))
}

func apiUpdateEmployerByID(c *gin.Context) {
	id, err := validationString(c.Param("id"))
	if err != nil {
		responseError(c, nil, err, 400)
		return
	}

	var employerPost MongoEmployerPost
	if err := c.BindJSON(&employerPost); err != nil {
		responseError(c, nil, err, 400)
		return
	}

	err = dbUpdateEmployee(id, employerPost.Name, employerPost.Gender, employerPost.DoB)
	if err != nil {
		responseError(c, nil, err, 404)
		return
	}
	delCache(id)

	msg := fmt.Sprintf("MongoEmployer %s was updated:\nName: %s\nGender: %s\nDoB: %s\n",
		id, employerPost.Name, employerPost.Gender, employerPost.DoB)
	responseOK(c, id, msg)
}

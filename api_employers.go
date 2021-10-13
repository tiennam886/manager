package manager

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"math"
	"net/http"
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
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"success":   false,
			"data":      nil,
			"total":     0,
			"page":      0,
			"last_page": 0,
			"limit":     0,
			"message":   err.Error(),
		})
		return
	}

	last := math.Ceil(float64(total / int64(limit)))
	if last < 1 && total > 0 {
		last = 1
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      employers,
		"total":     total,
		"page":      page,
		"last_page": last,
		"limit":     limit,
		"message":   "Get all employers successfully",
	})
	return
}

func PostEmployer(c *gin.Context) {
	var employer *Employer
	var err error

	if err := c.BindJSON(&employer); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data":    employer,
			"message": err.Error(),
		})
		return
	}

	employer.Name, employer.Gender, employer.DoB, err = validationAddEmployer(employer.Name, string(rune(employer.Gender)), employer.DoB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"data":    employer,
			"message": err.Error(),
		})
		return
	}

	//insert the newly created object into mongodb
	err = dbAddEmployer(employer.Name, employer.Gender, employer.DoB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"data":    employer,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    employer,
		"message": "Employer was created",
	})
	return

}

func DelEmployerByID(c *gin.Context) {
	id, err := validationString(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"data":    id,
			"message": err.Error()})
		return
	}

	err = dbDeleteEmployer(id)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"data":    id,
			"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusNoContent, gin.H{
		"success": true,
		"data":    id,
		"message": fmt.Sprintf("Hash with ID %s was deleted", id)})
	return
}

func UpdateEmployerByID(c *gin.Context) {
	var employer *Employer

	id, err := validationString(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"data":    id,
			"message": err.Error()})
		return
	}

	if err := c.BindJSON(&employer); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data":    employer,
			"message": err.Error(),
		})
		return
	}

	employer.Name, employer.Gender, employer.DoB, err = validationAddEmployer(employer.Name, string(rune(employer.Gender)), employer.DoB)

	err = dbUpdateEmployer(id, employer.Name, employer.Gender, employer.DoB)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"data":    employer,
			"message": err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"success": true,
		"data":    employer,
		"message": "Update successfully",
	})

	return
}

func response(success bool, data Employer, msg string) map[string]interface{} {

	res := bson.M{
		"success": success,
		"data":    data,
		"message": msg,
	}

	return res
}

package manager

import (
	"fmt"
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

	employers, total, err := employerMongo.ShowAll(page, limit)

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

	if err := c.BindJSON(&employer); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data":    employer,
			"message": err.Error(),
		})
		return
	}

	//insert the newly created object into mongodb
	err := employerMongo.AddEmployer(employer.Name, employer.Gender, employer.DoB)
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
	var id = c.Param("id")

	err := employerMongo.DeleteEmployer(id)
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
	var id = c.Param("id")

	if err := c.BindJSON(&employer); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data":    employer,
			"message": err.Error(),
		})
		return
	}

	err := employerMongo.UpdateEmployer(id, employer.Name, employer.Gender, employer.DoB)
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

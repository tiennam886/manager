package manager

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var serverHost = "localhost"
var serverPort = "8080"

func serverMode() error {
	router := gin.Default()

	router.GET("/employee", GetEmployers)
	router.GET("/team", GetTeams)
	router.GET("/team/:id", GetAllMemberInTeam)

	router.POST("/employee", PostEmployer)
	router.POST("/team", PostTeam)

	router.DELETE("/employee/:id", DelEmployerByID)
	router.DELETE("/team/:id", DelTeamByID)
	router.DELETE("/team/:id/employee/:mid", DelMemberInTeamByID)

	router.PATCH("/employee/:id", UpdateEmployerByID)
	router.PATCH("/team/:id/employee/:mid", AddMemberToTeamByID)
	router.PATCH("/team/:id", ChangeTeamName)

	var addr string
	addr = fmt.Sprintf("%s:%s", serverHost, serverPort)
	err := router.Run(addr)
	if err != nil {
		return err
	}

	return nil
}

func responseAllNotFound(c *gin.Context, err error) {
	c.IndentedJSON(http.StatusNotFound, gin.H{
		"success":   false,
		"data":      nil,
		"total":     0,
		"page":      0,
		"last_page": 0,
		"limit":     0,
		"message":   err.Error(),
	})
}

func responseAllEmployeeOK(c *gin.Context, data []Employer, total int64, page int, last float64, limit int) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      data,
		"total":     total,
		"page":      page,
		"last_page": last,
		"limit":     limit,
		"message":   "Get All Employer Successfully\n",
	})
}

func responseAllTeamOK(c *gin.Context, data []Teams, total int64, page int, last float64, limit int) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"success":   true,
		"data":      data,
		"total":     total,
		"page":      page,
		"last_page": last,
		"limit":     limit,
		"message":   "Get All Employer Successfully\n",
	})
}

func responseBadRequest(c *gin.Context, id string, err error) {
	c.IndentedJSON(http.StatusBadRequest, gin.H{
		"success": false,
		"data":    id,
		"message": err.Error(),
	})
}

func responseInternalServer(c *gin.Context, id string, err error) {
	c.IndentedJSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"data":    id,
		"message": err.Error(),
	})
}

func responseOK(c *gin.Context, id string, msg string) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"success": true,
		"data":    id,
		"message": msg,
	})
}

func responseEmployerOK(c *gin.Context, employer *EmployerPost, msg string) {
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    employer,
		"message": msg,
	})
}

func responseTeamCreated(c *gin.Context, team *Teams, msg string) {
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    team,
		"message": msg,
	})
}

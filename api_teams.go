package manager

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"strconv"
)

func GetTeams(c *gin.Context) {

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10
	}

	teams, total, err := teamMongo.ShowAllTeam(page, limit)

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
		"data":      teams,
		"total":     total,
		"page":      page,
		"last_page": last,
		"limit":     limit,
	})

	return
}

func GetAllMemberInTeam(c *gin.Context) {

	id := c.Param("id")

	team, err := teamMongo.ShowAllTeamMember(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"success": false,
			"date":    nil,
			"message": err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"success": true,
		"data":    team,
		"message": "All member in team has showed",
	})

	return

}

func PostTeam(c *gin.Context) {
	var team *Teams

	if err := c.BindJSON(&team); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data":    team,
			"message": err.Error(),
		})
		return
	}

	//insert the newly created object into mongodb
	err := teamMongo.AddTeam(team.Team)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"data":    team,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    team,
		"message": "Employer was created",
	})
	return
}

func DelTeamByID(c *gin.Context) {

	id := c.Param("id")
	err := teamMongo.DeleteTeamById(id)
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

func AddMemberToTeamByID(c *gin.Context) {

	if c.Param("id") == "" && c.Param("mid") == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data":    nil,
			"message": "request in format /teams/TEAM_ID/employers/EMPLOYER_ID",
		})
		return
	}

	var id = c.Param("id")
	var memberId = c.Param("mid")

	err1 := teamMongo.AddTeamMember(id, memberId)
	if err1 != nil {
		fmt.Println(err1)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Update failed "})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"success": true,
		"data":    id,
		"message": fmt.Sprintf("Add employer with id: %s in team with id: %s successfully", memberId, id),
	})
	return
}

func DelMemberInTeamByID(c *gin.Context) {

	if c.Param("id") == "" && c.Param("mid") == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data":    nil,
			"message": "request in format /teams/TEAM_ID/employers/EMPLOYER_ID",
		})
		return
	}

	var id = c.Param("id")
	var memberId = c.Param("mid")

	err := teamMongo.DelTeamMemberById(id, memberId)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"data":    nil,
			"message": err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"success": true,
		"data":    id,
		"message": fmt.Sprintf("Delete employer with id: %s in team with id: %s successfully", memberId, id),
	})
	return
}

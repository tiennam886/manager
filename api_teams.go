package manager

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func apiGetTeams(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 10
	}

	teams, total, err := dbGetAllTeams(page, limit)
	if err != nil {
		responseAllNotFound(c, err)
		return
	}

	last := math.Ceil(float64(total / int64(limit)))
	if last < 1 && total > 0 {
		last = 1
	}
	responseAllTeamOK(c, teams, total, page, last, limit)
}

func apiGetAllMemberInTeam(c *gin.Context) {
	id, err := validationString(c.Param("id"))
	if err != nil {
		responseBadRequest(c, c.Param("id"), err)
		return
	}

	team, err := dbShowAllMemberInTeam(id)
	if err != nil {
		responseInternalServer(c, id, err)
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"success": true,
		"data":    team,
		"message": "All member in team has showed\n",
	})
}

func apiPostTeam(c *gin.Context) {
	var (
		team *Teams
		err  error
	)

	if err := c.BindJSON(&team); err != nil {
		responseBadRequest(c, team.ID.Hex(), err)
		return
	}

	team.Team, err = validationString(team.Team)
	if err != nil {
		responseBadRequest(c, team.ID.Hex(), err)
		return
	}

	//insert the newly created object into mongodb
	err = dbAddTeam(team.Team)
	if err != nil {
		responseInternalServer(c, team.ID.Hex(), err)
		return
	}

	responseTeamCreated(c, team, "Team was created\n")
}

func apiDelTeamByID(c *gin.Context) {
	id, err := validationString(c.Param("id"))
	if err != nil {
		responseBadRequest(c, c.Param("id"), err)
		return
	}

	err = dbDeleteTeamById(id)
	if err != nil {
		responseInternalServer(c, id, err)
		return
	}

	responseOK(c, id, fmt.Sprintf("Team with ID %s was deleted\n", id))
}

func apiAddMemberToTeamByID(c *gin.Context) {
	id := c.Param("id")
	memberId := c.Param("mid")

	err = dbAddTeamMember(id, memberId)
	if err != nil {
		responseInternalServer(c, id, err)
		return
	}

	msg := fmt.Sprintf("Add employer with id: %s in team with id: %s successfully\n", memberId, id)
	responseOK(c, id, msg)
}

func apiDelMemberInTeamByID(c *gin.Context) {
	id := c.Param("id")
	mid := c.Param("mid")
	err = dbDelTeamMemberById(id, mid)
	if err != nil {
		responseInternalServer(c, id, err)
		return
	}

	msg := fmt.Sprintf("Delete employer with id: %s in team with id: %s successfully\n", mid, id)
	responseOK(c, id, msg)
}

func apiChangeTeamName(c *gin.Context) {
	var team *Teams

	id, err := validationString(c.Param("id"))
	if err != nil {
		responseBadRequest(c, c.Param("id"), err)
		return
	}

	if err = c.BindJSON(&team); err != nil {
		responseBadRequest(c, id, err)
		return
	}

	team.Team, err = validationString(team.Team)
	if err != nil {
		responseBadRequest(c, id, err)
		return
	}

	err = dbUpdateTeam(id, team.Team)
	if err != nil {
		responseInternalServer(c, team.ID.Hex(), err)
		return
	}
	responseOK(c, id, "Change Team name successfully\n")
}

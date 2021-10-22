package manager

import (
	json2 "encoding/json"
	"fmt"
	"math"
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

	teams, total, err := dbGetAllTeam(page, limit)
	if err != nil {
		responseError(c, nil, err, 404)
		return
	}

	last := math.Ceil(float64(int64(total) / int64(limit)))
	if last < 1 && total > 0 {
		last = 1
	}

	responseAllDataOK(c, teams, int64(total), page, last, limit)
}

func apiGetAllMemberInTeam(c *gin.Context) {
	var team MongoTeamMem

	id := c.Param("id")
	data, _ := getCache(id)
	err = json2.Unmarshal([]byte(data), &team)
	if data != "" && err == nil {
		responseOK(c, team, "All member in team has showed")
		return
	}

	resp, err := dbShowMemberInTeam(id)
	if err != nil {
		responseError(c, nil, err, 404)
		return
	}

	setCache(id, resp)

	responseOK(c, resp, "All member in team has showed")
}

func apiPostTeam(c *gin.Context) {
	var team *MongoTeam
	if err := c.BindJSON(&team); err != nil {
		responseError(c, nil, err, 400)
		return
	}

	err = dbAddTeam(team.Team)
	if err != nil {
		responseError(c, nil, err, 404)
		return
	}

	responseOK(c, team, "MySqlTeam was created\n")
}

func apiDelTeamByID(c *gin.Context) {
	id := c.Param("id")
	err = dbDelTeam(id)
	if err != nil {
		responseError(c, nil, err, 404)
		return
	}

	delCache(id)

	responseOK(c, id, fmt.Sprintf("MySqlTeam with ID %s was deleted\n", id))
}

func apiAddMemberToTeamByID(c *gin.Context) {
	id := c.Param("id")
	memberId := c.Param("mid")

	err = dbAddTeamMember(id, memberId)
	if err != nil {
		responseError(c, nil, err, 404)
		return
	}

	delCache(id)

	msg := fmt.Sprintf("Add employer with id: %s in team with id: %s successfully\n", memberId, id)
	responseOK(c, id, msg)
}

func apiDelMemberInTeamByID(c *gin.Context) {
	id := c.Param("id")
	mid := c.Param("mid")
	err = dbDelTeamMember(id, mid)
	if err != nil {
		responseError(c, nil, err, 404)
		return
	}

	delCache(id)

	msg := fmt.Sprintf("Delete employer with id: %s in team with id: %s successfully\n", mid, id)
	responseOK(c, id, msg)
}

func apiChangeTeamName(c *gin.Context) {
	var team *MongoTeam
	id := c.Param("id")

	if err = c.BindJSON(&team); err != nil {
		responseError(c, nil, err, 400)
		return
	}

	err = dbUpdateTeamName(id, team.Team)
	if err != nil {
		responseError(c, nil, err, 404)
		return
	}

	delCache(id)

	responseOK(c, id, "Change MySqlTeam name successfully\n")
}

package manager

import (
	"fmt"

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

	var addr string
	addr = fmt.Sprintf("%s:%s", serverHost, serverPort)
	err := router.Run(addr)
	if err != nil {
		return err
	}

	return nil
}

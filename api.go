package manager

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

var serverHost = "localhost"
var serverPort = "8080"

func serverMode() error {
	router := gin.Default()

	router.GET("/employers", GetEmployers)
	router.GET("/teams", GetTeams)
	router.GET("/teams/:id", GetAllMemberInTeam)

	router.POST("/employers", PostEmployer)
	router.POST("/teams", PostTeam)

	router.DELETE("/employers/:id", DelEmployerByID)
	router.DELETE("/teams/:id", DelTeamByID)
	router.DELETE("/teams/:id/employers/:mid", DelMemberInTeamByID)

	router.PATCH("/employers/:id", UpdateEmployerByID)
	router.PATCH("/teams/:id/employers/:mid", AddMemberToTeamByID)

	var addr string
	addr = fmt.Sprintf("%s:%s", serverHost, serverPort)
	err := router.Run(addr)
	if err != nil {
		return err
	}

	return nil
}

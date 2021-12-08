package service

import (
	"context"
	"github.com/tiennam886/manager/team/internal/model"
	"github.com/tiennam886/manager/team/internal/persistence"
)

func GetAllTeam(ctx context.Context, offset int, limit int) (teams []model.Team, err error) {
	teams, err = persistence.Teams().FindAll(ctx, offset, limit)
	return
}

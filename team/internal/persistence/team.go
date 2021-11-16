package persistence

import (
	"context"

	"github.com/tiennam886/manager/team/internal/model"
)

var teams TeamRepository

type TeamRepository interface {
	FindByUID(ctx context.Context, uid string) (model.Team, error)
	Save(ctx context.Context, team model.Team) error
	Update(ctx context.Context, uid string, team model.Team) error
	Remove(ctx context.Context, uid string) error
}

func Team() TeamRepository {
	if teams == nil {
		panic("persistence: employees not initiated")
	}
	return teams
}

func LoadTeamRepositoryWithMongoDB() (err error) {
	teams, err = newMongoTeamRepository()
	return err
}

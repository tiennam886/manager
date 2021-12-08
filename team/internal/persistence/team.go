package persistence

import (
	"context"
	"fmt"

	"github.com/tiennam886/manager/team/internal/model"
)

var (
	teams     TeamRepository
	teamDB    TeamRepository
	teamCache TeamRepository
)

type TeamRepository interface {
	FindAll(ctx context.Context, offset int, limit int) ([]model.Team, error)
	FindByUID(ctx context.Context, uid string) (model.Team, error)
	Save(ctx context.Context, team model.Team) error
	Update(ctx context.Context, uid string, team model.Team) error
	Remove(ctx context.Context, uid string) error

	AddAnEmployee(ctx context.Context, employeeId string, teamId string) error
	FindByTeamId(ctx context.Context, teamId string) ([]string, error)
	DeleteByTeamId(ctx context.Context, teamId string) error
	DeleteAnEmployee(ctx context.Context, employeeId string, teamId string) error
}

type TeamRepo struct {
	TeamDB    TeamRepository
	TeamCache TeamRepository
}

func (e TeamRepo) FindAll(ctx context.Context, offset int, limit int) ([]model.Team, error) {
	data, err := e.TeamDB.FindAll(ctx, offset, limit)
	return data, err
}

func (e TeamRepo) AddAnEmployee(ctx context.Context, employeeId string, teamId string) error {
	panic("implement me")
}

func (e TeamRepo) FindByTeamId(ctx context.Context, teamId string) ([]string, error) {
	panic("implement me")
}

func (e TeamRepo) DeleteByTeamId(ctx context.Context, teamId string) error {
	panic("implement me")
}

func (e TeamRepo) DeleteAnEmployee(ctx context.Context, employeeId string, teamId string) error {
	panic("implement me")
}

func (e TeamRepo) FindByUID(ctx context.Context, uid string) (model.Team, error) {
	data, err := e.TeamCache.FindByUID(ctx, uid)

	if err == nil {
		return data, err
	}

	data, err = e.TeamDB.FindByUID(ctx, uid)

	return data, e.TeamCache.Save(ctx, data)
}

func (e TeamRepo) Save(ctx context.Context, employee model.Team) error {
	return e.TeamDB.Save(ctx, employee)
}

func (e TeamRepo) Update(ctx context.Context, uid string, team model.Team) error {
	err := e.TeamDB.Update(ctx, uid, team)
	if err != nil {
		return err
	}

	return e.TeamCache.Remove(ctx, uid)
}

func (e TeamRepo) Remove(ctx context.Context, uid string) error {
	err := e.TeamDB.Remove(ctx, uid)
	if err != nil {
		return err
	}

	return e.TeamCache.Remove(ctx, uid)
}

func Teams() TeamRepository {
	if teams == nil {
		panic("persistence: teams not initiated")
	}

	return teams
}

func LoadTeamRepository(db string) (err error) {
	switch db {
	case "mongo":
		err = LoadTeamRepositoryWithMongoDB()
	case "postgres":
		err = LoadTeamRepositoryWithPostgresql()
	case "mysql":
		err = LoadTeamRepositoryWithMysql()
	default:
		err = fmt.Errorf("invalid database, choose mongo, postgres or mysql")
	}

	if err != nil {
		return err
	}

	err = LoadTeamRepositoryWithRedis()
	if err != nil {
		return
	}

	teams = &TeamRepo{
		TeamDB:    teamDB,
		TeamCache: teamCache,
	}
	return
}

func LoadTeamRepositoryWithMongoDB() (err error) {
	teamDB, err = newMongoTeamRepository()
	return
}

func LoadTeamRepositoryWithPostgresql() (err error) {
	teamDB, err = newPostgresqlTeamRepository()
	return
}

func LoadTeamRepositoryWithMysql() (err error) {
	teamDB, err = newMySqlTeamRepository()
	return
}

func LoadTeamRepositoryWithRedis() (err error) {
	teamCache, err = newRedisTeamRepository()
	return
}

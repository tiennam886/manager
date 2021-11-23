package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/tiennam886/manager/team/internal/config"
	"github.com/tiennam886/manager/team/internal/model"
)

type postgresqlTeamRepository struct {
	Database  *sql.DB
	TeamTable string
}

func newPostgresqlTeamRepository() (repo TeamRepository, err error) {
	postgresqlDB, err := sql.Open("postgres", config.Get().PostgresqlUrl)
	if err != nil {
		return
	}

	repo = &postgresqlTeamRepository{
		Database:  postgresqlDB,
		TeamTable: config.Get().TeamTable,
	}
	return
}

func (m postgresqlTeamRepository) FindByUID(ctx context.Context, uid string) (model.Team, error) {
	var team model.Team

	stmt := fmt.Sprintf("Select * from %s where uid = '%s';", m.TeamTable, uid)
	row, err := m.Database.Query(stmt)
	if err != nil {
		return team, err
	}
	row.Next()
	err = row.Scan(&team.UID, &team.Name, &team.Description)
	return team, err
}

func (m postgresqlTeamRepository) Save(ctx context.Context, team model.Team) error {
	stmt := fmt.Sprintf("insert into %s (uid, name, description) values($1, $2, $3)", m.TeamTable)
	_, err := m.Database.Exec(stmt, team.UID, team.Name, team.Description)
	return err
}

func (m postgresqlTeamRepository) Update(ctx context.Context, uid string, team model.Team) error {
	stmt := fmt.Sprintf("update %s set name=$1, description=$2 where uid=$3", m.TeamTable)
	_, err := m.Database.Exec(stmt, team.Name, team.Description, uid)
	return err
}

func (m postgresqlTeamRepository) Remove(ctx context.Context, uid string) error {
	stmt := fmt.Sprintf("delete from %s where uid=$1", m.TeamTable)
	_, err := m.Database.Exec(stmt, uid)
	return err
}

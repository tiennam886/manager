package persistence

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/tiennam886/manager/team/internal/config"
	"github.com/tiennam886/manager/team/internal/model"
)

type mysqlTeamRepository struct {
	Database  *sql.DB
	TeamTable string
}

func newMySqlTeamRepository() (repo TeamRepository, err error) {
	mySqlDB, err := sql.Open("mysql", config.Get().MysqlUrl)
	if err != nil {
		return
	}

	repo = &mysqlTeamRepository{
		Database:  mySqlDB,
		TeamTable: config.Get().TeamTable,
	}
	return
}

func (m mysqlTeamRepository) FindByUID(ctx context.Context, uid string) (model.Team, error) {
	var team model.Team

	stmt := fmt.Sprintf("Select * from %s where uid = %s", m.TeamTable, uid)
	row, err := m.Database.Query(stmt)
	if err != nil {
		return team, err
	}
	row.Next()
	err = row.Scan(&team.UID, &team.Name, &team.Description)
	return team, err
}

func (m mysqlTeamRepository) Save(ctx context.Context, team model.Team) error {
	stmt := fmt.Sprintf("insert into %s (uid, name, description) values(?, ?, ?);", m.TeamTable)
	_, err := m.Database.Exec(stmt, team.UID, team.Name, team.Description)
	return err
}

func (m mysqlTeamRepository) Update(ctx context.Context, uid string, team model.Team) error {
	stmt := fmt.Sprintf("update %s set name=?, description=? where uid=?", m.TeamTable)
	_, err := m.Database.Exec(stmt, team.Name, team.Description, uid)
	return err
}

func (m mysqlTeamRepository) Remove(ctx context.Context, uid string) error {
	stmt := fmt.Sprintf("delete from %s where uid=?", m.TeamTable)
	_, err := m.Database.Exec(stmt, uid)
	return err
}

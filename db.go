package manager

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"time"

	"database/sql"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	_ "github.com/go-sql-driver/mysql"
)

var (
	uri         = fmt.Sprintf("mongodb://%s:%s", conf.ServerHost, conf.MongoPort)
	database    = conf.MongoDatabase
	mySqlSource = fmt.Sprintf("%s@tcp(%s:%s)/%s", conf.MySqlUser, conf.MySqlHost, conf.MySqlPort, conf.MySqlDatabase)
	mySqlDB     *sql.DB
)

func connectCol(uri string, database string, col string) (*mongo.Collection, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	ctx := initCtx()
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client.Database(database).Collection(col), nil
}

func initCtx() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	return ctx
}

func connectMySql() (*sql.DB, error) {
	mySqlDB, err = sql.Open("mysql", mySqlSource)
	return mySqlDB, err
}

func dbAddEmployer(name string, gender string, date string) error {
	name, g, date, err := validationAddEmployer(name, gender, date)
	if err != nil {
		return err
	}

	if db == "mysql" {
		return dbMySqlAddEmployee(name, g, date)
	}

	return mongoAddEmployer(name, g, date)
}

func dbShowAllEmp(page int, limit int) (interface{}, int, error) {
	if db == "mysql" {
		data, total, err := dbMySqlShowAllEmployees(page, limit)
		return data, total, err
	}

	data, total, err := mongoShowAllEmployee(page, limit)
	return data, total, err
}

func dbGetEmployee(id string) (interface{}, error) {
	if db == "mysql" {
		mId, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}
		employerPost, err := dbMySqlGetEmployeeByID(mId)
		return employerPost, err
	}

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	employerPost, err := mongoFindEmployeeID(objId)
	return employerPost, err
}

func dbDelEmployee(id string) error {
	if db == "mysql" {
		mId, err := strconv.Atoi(id)
		if err != nil {
			return err
		}

		err = dbMySqlDelEmployeeByID(mId)
		return err
	}

	err = mongoDeleteEmployer(id)
	return err
}

func dbUpdateEmployee(id string, name string, gender string, dob string) error {
	name, g, dob, err := validationAddEmployer(name, gender, dob)
	if err != nil {
		return err
	}

	if db == "mysql" {
		mId, err := strconv.Atoi(id)
		if err != nil {
			return err
		}

		err = dbMySqlUpdateEmployee(mId, name, g, dob)
		return err
	}

	err = mongoUpdateEmployer(id, name, g, dob)
	return err
}

// MySqlTeam
func dbAddTeam(name string) error {
	name, err = validationString(name)
	if err != nil {
		return err
	}

	if db == "mysql" {
		return dbMySqlAddTeam(name)
	}

	return mongoAddTeam(name)
}

func dbGetAllTeam(page int, limit int) (interface{}, int, error) {
	if db == "mysql" {
		data, total, err := dbMySqlShowAllTeams(page, limit)
		return data, total, err
	}
	data, total, err := mongoGetAllTeams(page, limit)
	return data, total, err
}

func dbDelTeam(id string) error {
	if db == "mysql" {
		teamId, err := strconv.Atoi(id)
		if err != nil {
			return err
		}
		return dbMySqlDelTeamByID(teamId)
	}
	return mongoDeleteTeamById(id)
}

func dbUpdateTeamName(id string, name string) error {
	name, err := validationString(name)
	if err != nil {
		return err
	}
	if db == "mysql" {
		teamId, err := strconv.Atoi(id)
		if err != nil {
			return err
		}
		return dbMySqlUpdateTeam(teamId, name)
	}
	return mongoUpdateTeam(id, name)
}

func dbShowMemberInTeam(id string) (interface{}, error) {
	if db == "mysql" {
		teamId, err := strconv.Atoi(id)
		if err != nil {
			return nil, err
		}
		data, err := dbMySqlGetTeam(teamId)
		return data, err
	}
	data, err := mongoShowAllMemberInTeam(id)
	return data, err
}

func dbAddTeamMember(teamId string, memId string) error {
	if db == "mysql" {
		tId, mId, err := validationID(teamId, memId)
		if err != nil {
			return err
		}
		return dbMySqlAddTeamMember(tId, mId)
	}
	return mongoAddTeamMember(teamId, memId)
}

func dbDelTeamMember(teamId string, memId string) error {
	if db == "mysql" {
		tId, mId, err := validationID(teamId, memId)
		if err != nil {
			return err
		}
		return dbMySqlDelTeamMember(tId, mId)
	}
	return mongoDelTeamMemberById(teamId, memId)
}

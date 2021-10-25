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

var mySqlDB *sql.DB

func connectCol(col string) (*mongo.Collection, error) {
	uri := fmt.Sprintf("mongodb://%s:%s", conf.ServerHost, conf.MongoPort)

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

	return client.Database(conf.MongoDatabase).Collection(col), nil
}

func initCtx() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	return ctx
}

func connectMySql() (*sql.DB, error) {
	mySqlSource := fmt.Sprintf("%s@tcp(%s:%s)/%s", conf.MySqlUser, conf.MySqlHost, conf.MySqlPort, conf.MySqlDatabase)
	// root@(127.0.0.1:3306)/app
	mySqlDB, err = sql.Open("mysql", mySqlSource)
	return mySqlDB, err
}

func dbAddEmployer(name string, gender string, date string) (string, error) {
	name, g, date, err := validationAddEmployer(name, gender, date)
	if err != nil {
		return "", err
	}

	if db == "mysql" {
		id, err := dbMySqlAddEmployee(name, g, date)
		return strconv.FormatInt(id, 10), err
	}

	id, err := mongoAddEmployer(name, g, date)
	id = id.(primitive.ObjectID).Hex()
	return id.(string), err
}

func dbShowAllEmp(page int, limit int) (interface{}, int, error) {
	if db == "mysql" {
		data, total, err := dbMySqlShowAllEmployees(page, limit)
		displayMySqlEmployees(data)
		return data, total, err
	}

	data, total, err := mongoShowAllEmployee(page, limit)
	displayMongoEmployees(data)
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
func dbAddTeam(name string) (string, error) {
	name, err = validationString(name)
	if err != nil {
		return "", err
	}

	if db == "mysql" {
		id, err := dbMySqlAddTeam(name)
		if err != nil {
			return "", err
		}
		return strconv.FormatInt(id, 10), err
	}
	id, err := mongoAddTeam(name)
	if err != nil {
		return "", err
	}
	id = id.(primitive.ObjectID).Hex()
	return id.(string), nil
}

func dbGetAllTeam(page int, limit int) (interface{}, int, error) {
	if db == "mysql" {
		data, total, err := dbMySqlShowAllTeams(page, limit)
		displayMySqlTeams(data)
		return data, total, err
	}

	data, total, err := mongoGetAllTeams(page, limit)
	displayMongoTeams(data)
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
		if err != nil {
			return nil, err
		}

		displayMySqlMembers(data)
		return data, err
	}

	data, err := mongoShowAllMemberInTeam(id)
	if err != nil {
		return nil, err
	}

	displayMongoTeamMembers(data)
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

//display

func displayMySqlEmployees(data interface{}) {
	employees := data.([]MySqlEmployee)

	fmt.Printf("ID\tNAME\t\tGENDER\tDOB\n")
	for i := range employees {
		fmt.Printf("%v\t%s\t%v\t%s\n",
			employees[i].ID, employees[i].Name, employees[i].Gender, employees[i].DoB)
	}
}

func displayMongoEmployees(data interface{}) {
	employees := data.([]MongoEmployerPost)

	fmt.Printf("ID\t\t\t\tNAME\t\tGENDER\tDOB\n")
	for i := range employees {
		fmt.Printf("%v\t%s\t%s\t%s\n",
			employees[i].ID.Hex(), employees[i].Name, employees[i].Gender, employees[i].DoB)
	}
}

func displayMySqlTeams(data interface{}) {
	teams := data.([]MySqlTeam)
	fmt.Printf("ID\tNAME\t\n")
	for i := range teams {
		fmt.Printf("%v\t%s\n", teams[i].ID, teams[i].Name)
	}
}

func displayMongoTeams(data interface{}) {
	teams := data.([]MongoTeam)
	fmt.Printf("ID\t\t\t\tNAME\t\n")
	for i := range teams {
		fmt.Printf("%v\t%s\n", teams[i].ID.Hex(), teams[i].Name)
	}
}

func displayMongoTeamMembers(data interface{}) {
	teamMem := data.(MongoTeamMem)
	fmt.Printf("List Employers in: %s with id: %s\n\n", teamMem.Team, teamMem.ID.Hex())
	fmt.Printf("ID\t\t\t\tNAME\t\tGENDER\tDOB\n")
	for i := range teamMem.Member {
		fmt.Printf("%v\t%s\t%s\t%s\n",
			teamMem.Member[i].ID.Hex(), teamMem.Member[i].Name,
			convertNumToGender(teamMem.Member[i].Gender), teamMem.Member[i].DoB)
	}
}

func displayMySqlMembers(data interface{}) {
	teamMem := data.(MySqlTeamMem)
	fmt.Printf("List Employers in: %s with id: %v\n\n", teamMem.Name, teamMem.ID)
	displayMySqlEmployees(teamMem.Members)
}

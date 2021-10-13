package manager

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var teamCollection = "teams"

type Teams struct {
	ID     primitive.ObjectID   `bson:"_id"`
	Team   string               `bson:"team"`
	Member []primitive.ObjectID `bson:"member"`
}

type TeamMem struct {
	ID     primitive.ObjectID `bson:"_id"`
	Team   string             `bson:"team"`
	Member []Employer         `bson:"employers"`
}

func dbAddTeam(name string) error {
	var ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)

	team := bson.M{
		"team":   name,
		"member": []string{},
	}

	tCol, err := connectCol(uri, database, teamCollection)
	if err != nil {
		return err
	}

	resp, err := tCol.InsertOne(ctx, team)
	if err != nil {
		return err
	}

	fmt.Printf("Insert team with name %s to DB successfully with ID: %s\n", name, resp.InsertedID)

	return nil

}

func dbGetAllTeams(page int, limit int) ([]Teams, int64, error) {
	var ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)

	filter := bson.M{}

	findOptions := options.Find()
	findOptions.SetSkip((int64(page) - 1) * int64(limit))
	findOptions.SetLimit(int64(limit))

	tCol, err := connectCol(uri, database, teamCollection)
	if err != nil {
		return nil, 0, err
	}

	total, err := tCol.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	cursor, err := tCol.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}

	var teams []Teams
	err = cursor.All(ctx, &teams)
	if err != nil {
		return nil, 0, err
	}

	return teams, total, nil
}

func dbShowAllMemberInTeam(id string) (TeamMem, error) {
	var resp []TeamMem

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return resp[0], fmt.Errorf("TEAM_ID %s was invalid\n", id)
	}

	pipeline := []bson.M{
		{"$match": bson.M{"_id": objId}},
		{"$lookup": bson.M{
			"from":         "employers",
			"localField":   "member",
			"foreignField": "_id",
			"as":           "employers"},
		},
	}

	tCol, err := connectCol(uri, database, teamCollection)
	if err != nil {
		return resp[0], err
	}

	cursor, err := tCol.Aggregate(ctx, pipeline)
	if err != nil {
		return resp[0], err
	}
	err = cursor.All(ctx, &resp)
	if err != nil {
		return resp[0], err
	}
	return resp[0], nil
}

func dbAddTeamMember(id string, newMemberId string) error {
	objId, memId, err := validationObjectID(id, newMemberId)
	if err != nil {
		return err
	}

	err = dbFindEmployeeID(memId)
	if err != nil {
		return err
	}

	team, err := dbFindTeamID(objId)
	if err != nil {
		return err
	}

	for _, id := range team.Member {
		if memId == id {
			return fmt.Errorf("Member with id: %s has already been in %s\n", memId, team.Team)
		}
	}

	err = dbUpdateTeamMember(objId, memId, "$push")
	if err != nil {
		return err
	}
	fmt.Printf("User with id: %s was added to Team with id: %s \n", newMemberId, id)
	return nil
}

func dbDelTeamMemberById(id string, delMemberId string) error {
	objId, memId, err := validationObjectID(id, delMemberId)
	if err != nil {
		return err
	}

	err = dbUpdateTeamMember(objId, memId, "$pull")
	if err != nil {
		return err
	}

	fmt.Printf("User with id: %s was deleted from Team with id: %s \n", delMemberId, id)
	return nil
}

func dbUpdateTeamMember(objId primitive.ObjectID, memId primitive.ObjectID, method string) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	filter := bson.M{"_id": objId}
	update := bson.M{
		method: bson.M{"member": memId},
	}

	tCol, err := connectCol(uri, database, teamCollection)
	if err != nil {
		return err
	}

	_, err = tCol.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func dbFindTeamID(objId primitive.ObjectID) (Teams, error) {
	var team Teams

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	tCol, err := connectCol(uri, database, teamCollection)
	if err != nil {
		return team, err
	}

	err = tCol.FindOne(ctx, bson.M{"_id": objId}).Decode(&team)
	if err != nil {
		return team, fmt.Errorf("Team with id:  %s was not found\n", objId)
	}
	return team, nil
}

func dbDeleteTeamById(id string) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("TEAM_ID %s was invalid\n", id)
	}

	_, err = dbFindTeamID(objId)
	if err != nil {
		return err
	}

	tCol, err := connectCol(uri, database, teamCollection)
	if err != nil {
		return err
	}

	_, err = tCol.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return err
	}

	fmt.Printf("Team %s was deleted\n", id)
	return nil

}

func dbUpdateTeam(id string, name string) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("Employer_ID %s was invalid\n", id)
	}

	team, err := dbFindTeamID(objId)
	if err != nil {
		return err
	}

	newTeam := Teams{
		ID:     objId,
		Team:   name,
		Member: team.Member,
	}
	update := bson.M{
		"$set": newTeam,
	}

	tCol, err := connectCol(uri, database, teamCollection)
	if err != nil {
		return err
	}
	_, err = tCol.UpdateOne(ctx, bson.M{"_id": objId}, update)
	if err != nil {
		return err
	}
	fmt.Printf("Team %s was updated with new name: %s", id, name)
	return nil
}

package manager

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoTeam struct {
	ID     primitive.ObjectID   `bson:"_id" json:"id"`
	Team   string               `bson:"team" json:"team"`
	Member []primitive.ObjectID `bson:"member" json:"member"`
}

type MongoTeamMem struct {
	ID     primitive.ObjectID `bson:"_id" json:"id"`
	Team   string             `bson:"team" json:"team"`
	Member []MongoEmployer    `bson:"employers" json:"employers"`
}

func mongoAddTeam(name string) error {
	ctx := initCtx()

	team := bson.M{
		"team":   name,
		"member": []string{},
	}
	_, err := teamCol.InsertOne(ctx, team)
	return err
}

func mongoGetAllTeams(page int, limit int) (interface{}, int, error) {
	ctx := initCtx()

	filter := bson.M{}

	findOptions := options.Find()
	findOptions.SetSkip((int64(page) - 1) * int64(limit))
	findOptions.SetLimit(int64(limit))

	total, err := teamCol.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	cursor, err := teamCol.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}

	var teams []MongoTeam
	err = cursor.All(ctx, &teams)
	if err != nil {
		return nil, 0, err
	}

	return teams, int(total), nil
}

func mongoShowAllMemberInTeam(id string) (interface{}, error) {
	var resp []MongoTeamMem

	ctx := initCtx()

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

	cursor, err := teamCol.Aggregate(ctx, pipeline)
	if err != nil {
		return resp[0], err
	}
	err = cursor.All(ctx, &resp)
	if err != nil {
		return resp[0], err
	}
	return resp[0], nil
}

func mongoAddTeamMember(id string, newMemberId string) error {
	objId, memId, err := validationObjectID(id, newMemberId)
	if err != nil {
		return err
	}

	_, err = mongoFindEmployeeID(memId)
	if err != nil {
		return err
	}

	team, err := mongoFindTeamID(objId)
	if err != nil {
		return err
	}
	teamTransforms := team.(MongoTeam)
	for _, id := range teamTransforms.Member {
		if memId == id {
			return fmt.Errorf("Member with id: %s has already been in %s\n", memId, teamTransforms.Team)
		}
	}

	err = mongoUpdateTeamMember(objId, memId, "$push")
	if err != nil {
		return err
	}
	fmt.Printf("User with id: %s was added to MySqlTeam with id: %s \n", newMemberId, id)
	return nil
}

func mongoDelTeamMemberById(id string, delMemberId string) error {
	objId, memId, err := validationObjectID(id, delMemberId)
	if err != nil {
		return err
	}

	err = mongoUpdateTeamMember(objId, memId, "$pull")
	if err != nil {
		return err
	}

	fmt.Printf("User with id: %s was deleted from MySqlTeam with id: %s \n", delMemberId, id)
	return nil
}

func mongoUpdateTeamMember(objId primitive.ObjectID, memId primitive.ObjectID, method string) error {
	ctx := initCtx()

	filter := bson.M{"_id": objId}
	update := bson.M{
		method: bson.M{"member": memId},
	}

	_, err = teamCol.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func mongoFindTeamID(objId primitive.ObjectID) (interface{}, error) {
	var team MongoTeam

	ctx := initCtx()

	err = teamCol.FindOne(ctx, bson.M{"_id": objId}).Decode(&team)
	if err != nil {
		return team, fmt.Errorf("MySqlTeam with id:  %s was not found\n", objId)
	}
	return team, nil
}

func mongoDeleteTeamById(id string) error {
	ctx := initCtx()

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("TEAM_ID %s was invalid\n", id)
	}

	_, err = mongoFindTeamID(objId)
	if err != nil {
		return err
	}

	_, err = teamCol.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return err
	}

	fmt.Printf("MySqlTeam %s was deleted\n", id)
	return nil

}

func mongoUpdateTeam(id string, name string) error {
	ctx := initCtx()

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("Employer_ID %s was invalid\n", id)
	}

	name, err = validationString(name)
	if err != nil {
		return err
	}

	team, err := mongoFindTeamID(objId)
	if err != nil {
		return err
	}
	teamTransforms := team.(MongoTeam)
	newTeam := MongoTeam{
		ID:     objId,
		Team:   name,
		Member: teamTransforms.Member,
	}
	update := bson.M{
		"$set": newTeam,
	}
	_, err = teamCol.UpdateOne(ctx, bson.M{"_id": objId}, update)
	if err != nil {
		return err
	}
	fmt.Printf("MySqlTeam %s was updated with new name: %s\n", id, name)
	return nil
}

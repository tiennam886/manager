package manager

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

type TeamMongo struct {
	collection *mongo.Collection
}

func (h *TeamMongo) InitTeamRepo() error {
	db, err := ConnectDB(uri, database)
	if err != nil {
		return err
	}

	h.collection = db.Collection(teamCollection)
	return nil
}

func (h *TeamMongo) AddTeam(name string) error {
	var ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)

	team := bson.M{
		"team":   name,
		"member": []string{},
	}

	resp, err := h.collection.InsertOne(ctx, team)
	if err != nil {
		return err
	}

	fmt.Printf("Insert team with name %s to DB successfully with ID: %s\n", name, resp.InsertedID)

	return nil

}

func (h *TeamMongo) ShowAllTeam(page int, limit int) ([]Teams, int64, error) {
	var ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)

	filter := bson.M{}

	findOptions := options.Find()
	findOptions.SetSkip((int64(page) - 1) * int64(limit))
	findOptions.SetLimit(int64(limit))

	total, err := h.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	cursor, err := h.collection.Find(ctx, filter, findOptions)
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

func (h *TeamMongo) ShowAllTeamMember(id string) ([]TeamMem, error) {
	var ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)

	objId, err1 := primitive.ObjectIDFromHex(id)
	if err1 != nil {
		return nil, fmt.Errorf("TEAM_ID %s was invalid\n", id)
	}

	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"_id": objId,
			}},
		bson.M{
			"$lookup": bson.M{
				"from":         "employers",
				"localField":   "member",
				"foreignField": "_id",
				"as":           "employers",
			},
		},
	}

	cursor, err := h.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var resp []TeamMem
	err = cursor.All(ctx, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil

}

func (h *TeamMongo) AddTeamMember(id string, newMemberId string) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	objId, err1 := primitive.ObjectIDFromHex(id)
	if err1 != nil {
		return fmt.Errorf("TEAM_ID %s was invalid\n", id)
	}

	memId, err2 := primitive.ObjectIDFromHex(newMemberId)
	if err2 != nil {
		return fmt.Errorf("Employer_ID %s was invalid", newMemberId)
	}

	var employer Employer
	err := h.collection.FindOne(ctx, bson.M{"_id": memId}).Decode(&employer)
	if err != nil {
		return fmt.Errorf("Employers with ID %s does not exist.\n", newMemberId)
	}

	var team Teams
	err = h.collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&team)
	if err != nil {
		return fmt.Errorf("Team with id:  %s was not found\n", id)
	}

	for _, id := range team.Member {
		if memId == id {
			return fmt.Errorf("Member with id: %s has already been in   %s\n", memId, team.Team)
		}
	}

	filter := bson.M{"_id": objId}

	update := bson.M{
		"$push": bson.M{"member": memId},
	}

	_, err = h.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	fmt.Printf("User with id: %s was added to Team with id: %s \n", newMemberId, id)
	return nil
}

func (h *TeamMongo) DelTeamMemberById(id string, delMemberId string) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	objId, err1 := primitive.ObjectIDFromHex(id)
	if err1 != nil {
		return fmt.Errorf("TEAM_ID %s was invalid\n", id)
	}

	memId, err2 := primitive.ObjectIDFromHex(delMemberId)
	if err2 != nil {
		return fmt.Errorf("Employer_ID %s was invalid", delMemberId)
	}

	filter := bson.M{"_id": objId}

	update := bson.M{
		"$pull": bson.M{"member": memId},
	}

	_, err := h.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("Team with id:  %s was not found", id)
	}

	fmt.Printf("User with id: %s was deleted from Team with id: %s \n", delMemberId, id)
	return nil
}

func (h *TeamMongo) DeleteTeamById(id string) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	objId, err1 := primitive.ObjectIDFromHex(id)
	if err1 != nil {
		return fmt.Errorf("TEAM_ID %s was invalid\n", id)
	}

	filter := bson.M{"_id": objId}

	_, err := h.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("Team with id:  %s was not found", id)
	}

	fmt.Printf("Team %s was deleted\n", id)
	return nil

}

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

var employerCollection = "employers"

type Employer struct {
	ID     primitive.ObjectID `bson:"_id"`
	Name   string             `bson:"name"`
	Gender int                `bson:"gender"`
	DoB    string             `bson:"dob"`
}

type EmployerMongo struct {
	collection *mongo.Collection
}

func (h *EmployerMongo) InitEmployerRepo() error {
	db, err := ConnectDB(uri, database)
	if err != nil {
		return err
	}

	h.collection = db.Collection(employerCollection)
	return nil
}

func (h *EmployerMongo) AddEmployer(name string, gender int, date string) error {
	var ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)

	employer := bson.M{
		"name":   name,
		"gender": gender,
		"dob":    date,
	}

	resp, err := h.collection.InsertOne(ctx, employer)
	if err != nil {
		return err
	}

	fmt.Printf("Insert employer name: %s, gender: %v, DoB: %s to DB successfully with ID: %s\n", name, gender, date, resp.InsertedID)
	return nil
}

func (h *EmployerMongo) ShowAll(page int, limit int) ([]Employer, int64, error) {
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

	var employers []Employer
	err = cursor.All(ctx, &employers)
	if err != nil {
		return nil, 0, err
	}

	return employers, total, nil

}

func (h *EmployerMongo) UpdateEmployer(id string, newName string, newGender int, newBoB string) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	objId, err1 := primitive.ObjectIDFromHex(id)
	if err1 != nil {
		return fmt.Errorf("Employer_ID %s was invalid\n", id)
	}

	filter := bson.M{"_id": objId}

	newEmp := Employer{
		ID:     objId,
		Name:   newName,
		Gender: newGender,
		DoB:    newBoB,
	}
	update := bson.M{
		"$set": newEmp,
	}

	_, err := h.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	fmt.Printf("Employer %s was updated\nName: %s\nGender: %v\nDoB: %s\n", id, newName, newGender, newBoB)
	return nil

}

func (h *EmployerMongo) DeleteEmployer(id string) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	objId, err1 := primitive.ObjectIDFromHex(id)
	if err1 != nil {
		return fmt.Errorf("Employer_ID %s was invalid\n", id)
	}

	filter := bson.M{"_id": objId}

	_, err := h.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("Employer_ID %s does not exist\n", id)
	}

	fmt.Printf("Employer %s was deleted\n", id)
	return nil

}

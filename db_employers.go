package manager

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var employerCollection = "employers"

type Employer struct {
	ID     primitive.ObjectID `bson:"_id"`
	Name   string             `bson:"name"`
	Gender int                `bson:"gender"`
	DoB    string             `bson:"dob"`
}

func dbAddEmployer(name string, gender int, date string) error {
	var ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)

	employer := bson.M{
		"name":   name,
		"gender": gender,
		"dob":    date,
	}

	eCol, err := connectCol(uri, database, employerCollection)
	if err != nil {
		return err
	}

	resp, err := eCol.InsertOne(ctx, employer)
	if err != nil {
		return err
	}

	fmt.Printf("Insert employer name: %s, gender: %v, DoB: %s to DB successfully with ID: %s\n", name, gender, date, resp.InsertedID)
	return nil
}

func dbShowAllEmployee(page int, limit int) ([]Employer, int64, error) {
	var ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)

	filter := bson.M{}

	findOptions := options.Find()
	findOptions.SetSkip((int64(page) - 1) * int64(limit))
	findOptions.SetLimit(int64(limit))

	eCol, err := connectCol(uri, database, employerCollection)
	if err != nil {
		return nil, 0, err
	}

	total, err := eCol.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	cursor, err := eCol.Find(ctx, filter, findOptions)
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

func dbUpdateEmployer(id string, newName string, newGender int, newBoB string) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("Employer_ID %s was invalid\n", id)
	}

	newEmp := Employer{
		ID:     objId,
		Name:   newName,
		Gender: newGender,
		DoB:    newBoB,
	}
	update := bson.M{
		"$set": newEmp,
	}

	eCol, err := connectCol(uri, database, employerCollection)
	if err != nil {
		return err
	}
	_, err = eCol.UpdateOne(ctx, bson.M{"_id": objId}, update)
	if err != nil {
		return err
	}

	fmt.Printf("Employer %s was updated\nName: %s\nGender: %v\nDoB: %s\n", id, newName, newGender, newBoB)
	return nil

}

func dbDeleteEmployer(id string) error {
	var employer *Employer

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("Employer_ID %s was invalid\n", id)
	}

	filter := bson.M{"_id": objId}

	eCol, err := connectCol(uri, database, employerCollection)
	if err != nil {
		return err
	}

	err = eCol.FindOne(ctx, filter).Decode(&employer)
	if err != nil {
		return fmt.Errorf("Employer_ID %s does not exist\n", id)
	}

	_, err = eCol.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	fmt.Printf("Employer %s was deleted\n", id)
	return nil

}

func dbFindEmployeeID(id primitive.ObjectID) error {
	var employer Employer

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	eCol, err := connectCol(uri, database, employerCollection)
	if err != nil {
		return err
	}
	err = eCol.FindOne(ctx, bson.M{"_id": id}).Decode(&employer)
	if err != nil {
		return fmt.Errorf("Employers with ID %s does not exist.\n", id)
	}
	return nil
}

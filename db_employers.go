package manager

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var employerCollection = "employers"

type Employer struct {
	ID     primitive.ObjectID `bson:"_id" json:"id"`
	Name   string             `bson:"name" json:"name"`
	Gender int                `bson:"gender" json:"gender"`
	DoB    string             `bson:"dob" json:"dob"`
}

type EmployerPost struct {
	ID     primitive.ObjectID `bson:"_id" json:"id"`
	Name   string             `bson:"name" json:"name"`
	Gender string             `bson:"gender" json:"gender"`
	DoB    string             `bson:"dob" json:"dob"`
}

func dbAddEmployer(name string, gender int, date string) error {
	ctx := initCtx()

	employer := bson.M{
		"name":   name,
		"gender": gender,
		"dob":    date,
	}

	resp, err := employeeCol.InsertOne(ctx, employer)
	if err != nil {
		return err
	}

	fmt.Printf("Insert employer name: %s, gender: %s, DoB: %s to DB successfully with ID: %s\n",
		name, convertNumToGender(gender), date, resp.InsertedID)
	return nil
}

func dbShowAllEmployee(page int, limit int) ([]Employer, int64, error) {
	ctx := initCtx()

	filter := bson.M{}

	findOptions := options.Find()
	findOptions.SetSkip((int64(page) - 1) * int64(limit))
	findOptions.SetLimit(int64(limit))

	total, err := employeeCol.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	cursor, err := employeeCol.Find(ctx, filter, findOptions)
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
	ctx := initCtx()

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

	_, err = employeeCol.UpdateOne(ctx, bson.M{"_id": objId}, update)
	if err != nil {
		return err
	}

	fmt.Printf("Employer %s was updated:\nName: %s\nGender: %s\nDoB: %s\n",
		id, newName, convertNumToGender(newGender), newBoB)
	return nil
}

func dbDeleteEmployer(id string) error {
	var employer *Employer

	ctx := initCtx()

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("Employer_ID %s was invalid\n", id)
	}

	filter := bson.M{"_id": objId}

	err = employeeCol.FindOne(ctx, filter).Decode(&employer)
	if err != nil {
		return fmt.Errorf("Employer_ID %s does not exist\n", id)
	}

	_, err = employeeCol.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	fmt.Printf("Employer %s was deleted\n", id)
	return nil

}

func dbFindEmployeeID(id primitive.ObjectID) error {
	var employer Employer

	ctx := initCtx()

	err := employeeCol.FindOne(ctx, bson.M{"_id": id}).Decode(&employer)
	if err != nil {
		return fmt.Errorf("Employers with ID %s does not exist.\n", id)
	}
	return nil
}

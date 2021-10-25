package manager

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	employer  MongoEmployer
	employers []MongoEmployerPost
)

type MongoEmployer struct {
	ID     primitive.ObjectID `bson:"_id" json:"id"`
	Name   string             `bson:"name" json:"name"`
	Gender int                `bson:"gender" json:"gender"`
	DoB    string             `bson:"dob" json:"dob"`
}

type MongoEmployerPost struct {
	ID     primitive.ObjectID `bson:"_id" json:"id"`
	Name   string             `bson:"name" json:"name"`
	Gender string             `bson:"gender" json:"gender"`
	DoB    string             `bson:"dob" json:"dob"`
}

func mongoAddEmployer(name string, gender int, date string) (interface{}, error) {
	ctx := initCtx()

	employer := bson.M{
		"name":   name,
		"gender": gender,
		"dob":    date,
	}
	resp, err := employeeCol.InsertOne(ctx, employer)
	return resp.InsertedID, err
}

func mongoShowAllEmployee(page int, limit int) (interface{}, int, error) {
	ctx := initCtx()

	filter := bson.M{}

	findOptions := options.Find()
	findOptions.SetSkip((int64(page) - 1) * int64(limit))
	findOptions.SetLimit(int64(limit))

	cursor, err := employeeCol.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}

	for cursor.Next(ctx) {
		err = cursor.Decode(&employer)
		if err != nil {
			return employers, 0, nil
		}
		e := MongoEmployerPost{
			ID:     employer.ID,
			Name:   employer.Name,
			Gender: convertNumToGender(employer.Gender),
			DoB:    employer.DoB,
		}
		employers = append(employers, e)
	}

	return employers, len(employers), nil
}

func mongoUpdateEmployer(id string, newName string, newGender int, newDoB string) error {
	ctx := initCtx()

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("Employer_ID %s was invalid\n", id)
	}

	newEmp := MongoEmployer{
		ID:     objId,
		Name:   newName,
		Gender: newGender,
		DoB:    newDoB,
	}
	update := bson.M{
		"$set": newEmp,
	}

	_, err = employeeCol.UpdateOne(ctx, bson.M{"_id": objId}, update)
	if err != nil {
		return err
	}

	fmt.Printf("MongoEmployer %s was updated:\nName: %s\nGender: %s\nDoB: %s\n",
		id, newName, convertNumToGender(newGender), newDoB)
	return nil
}

func mongoDeleteEmployer(id string) error {
	ctx := initCtx()

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("Employer_ID %s was invalid\n", id)
	}

	filter := bson.M{"_id": objId}

	_, err = mongoFindEmployeeID(objId)
	if err != nil {
		return err
	}

	_, err = employeeCol.DeleteOne(ctx, filter)
	return err
}

func mongoFindEmployeeID(id primitive.ObjectID) (interface{}, error) {
	ctx := initCtx()

	var employerPost MongoEmployerPost
	err := employeeCol.FindOne(ctx, bson.M{"_id": id}).Decode(&employer)
	if err != nil {
		return employerPost, fmt.Errorf("Employers with ID %s does not exist.\n", id)
	}
	employerPost = MongoEmployerPost{
		ID:     employer.ID,
		Name:   employer.Name,
		Gender: convertNumToGender(employer.Gender),
		DoB:    employer.DoB,
	}
	return employerPost, nil
}

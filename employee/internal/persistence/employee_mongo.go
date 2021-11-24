package persistence

import (
	"context"
	"time"

	"github.com/tiennam886/manager/employee/internal/config"
	"github.com/tiennam886/manager/employee/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoEmployeeRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func (repo *mongoEmployeeRepository) FindAll(ctx context.Context, offset int, limit int) ([]model.EmployeePost, error) {
	filter := bson.M{}

	findOptions := options.Find()
	findOptions.SetSkip((int64(offset) - 1) * int64(limit))
	findOptions.SetLimit(int64(limit))

	cursor, err := repo.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return []model.EmployeePost{}, err
	}

	var employees []model.EmployeePost
	for cursor.Next(ctx) {
		var employee model.Employee
		err = cursor.Decode(&employee)
		if err != nil {
			return employees, err
		}
		employees = append(employees, employee.ToEmployeePost())
	}
	return employees, nil
}

func newMongoEmployeeRepository() (repo EmployeeRepository, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Get().MongoDbUrl))
	if err != nil {
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return
	}

	repo = &mongoEmployeeRepository{
		client:     client,
		collection: client.Database(config.Get().Database).Collection(config.Get().Collection),
	}
	return repo, nil
}

func (repo *mongoEmployeeRepository) FindByUID(ctx context.Context, uid string) (model.EmployeePost, error) {
	result := repo.collection.FindOne(ctx, bson.M{"uid": uid})

	var staff model.Employee
	if err := result.Decode(&staff); err != nil {
		return model.EmployeePost{}, err
	}

	return staff.ToEmployeePost(), nil
}

func (repo *mongoEmployeeRepository) Save(ctx context.Context, staff model.Employee) error {
	_, err := repo.collection.InsertOne(ctx, staff)
	return err
}

func (repo *mongoEmployeeRepository) Update(ctx context.Context, uid string, staff model.Employee) error {
	_, err := repo.collection.UpdateOne(ctx, bson.M{"uid": uid}, toStaffDocument(staff))
	return err
}

func (repo *mongoEmployeeRepository) Remove(ctx context.Context, uid string) error {
	_, err := repo.collection.DeleteOne(ctx, bson.M{"uid": uid})
	return err
}

type EmployeeDocument struct {
	ID     primitive.ObjectID `bson:"_id"`
	UID    string             `bson:"uid"`
	Name   string             `bson:"name"`
	DOB    string             `bson:"dob"`
	Gender int                `bson:"gender"`
}

func toStaffDocument(s model.Employee) EmployeeDocument {
	return EmployeeDocument{
		ID:     primitive.NewObjectID(),
		UID:    s.UID,
		Name:   s.Name,
		DOB:    s.DOB,
		Gender: s.Gender,
	}
}

func (s EmployeeDocument) ToModel() model.Employee {
	return model.Employee{
		UID:    s.UID,
		Name:   s.Name,
		DOB:    s.DOB,
		Gender: s.Gender,
	}
}

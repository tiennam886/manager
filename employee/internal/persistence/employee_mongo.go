package persistence

import (
	"context"
	"github.com/tiennam886/manager/employee/internal/config"
	"github.com/tiennam886/manager/employee/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type mongoStaffRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func newMongoStaffRepository() (repo EmployeeRepository, err error) {
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

	repo = &mongoStaffRepository{
		client:     client,
		collection: client.Database(config.Get().Database).Collection(config.Get().Collection),
	}
	return repo, nil
}

func (repo *mongoStaffRepository) FindByUID(ctx context.Context, uid string) (model.Employee, error) {
	result := repo.collection.FindOne(ctx, bson.M{"uid": uid})

	var staff model.Employee
	if err := result.Decode(&staff); err != nil {
		return model.Employee{}, err
	}

	return staff, nil
}

func (repo *mongoStaffRepository) Save(ctx context.Context, staff model.Employee) error {
	_, err := repo.collection.InsertOne(ctx, staff)
	return err
}

func (repo *mongoStaffRepository) Update(ctx context.Context, uid string, staff model.Employee) error {
	_, err := repo.collection.UpdateOne(ctx, bson.M{"uid": uid}, toStaffDocument(staff))
	return err
}

func (repo *mongoStaffRepository) Remove(ctx context.Context, uid string) error {
	_, err := repo.collection.DeleteOne(ctx, bson.M{"uid": uid})
	return err
}

type EmployeeDocument struct {
	ID     primitive.ObjectID `bson:"_id"`
	UID    string             `bson:"uid"`
	Name   string             `bson:"name"`
	DOB    primitive.DateTime `bson:"dob"`
	Gender int                `bson:"gender"`
}

func toStaffDocument(s model.Employee) EmployeeDocument {
	return EmployeeDocument{
		ID:     primitive.NewObjectID(),
		UID:    s.UID,
		Name:   s.Name,
		DOB:    primitive.NewDateTimeFromTime(s.DOB),
		Gender: s.Gender,
	}
}

func (s EmployeeDocument) ToModel() model.Employee {
	return model.Employee{
		UID:    s.UID,
		Name:   s.Name,
		DOB:    s.DOB.Time(),
		Gender: s.Gender,
	}
}

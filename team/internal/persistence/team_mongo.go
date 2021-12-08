package persistence

import (
	"context"
	"time"

	"github.com/tiennam886/manager/team/internal/config"
	"github.com/tiennam886/manager/team/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoTeamRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
	members    *mongo.Collection
}

func (repo *mongoTeamRepository) FindAll(ctx context.Context, offset int, limit int) ([]model.Team, error) {
	filter := bson.M{}

	findOptions := options.Find()
	findOptions.SetSkip((int64(offset) - 1) * int64(limit))
	findOptions.SetLimit(int64(limit))

	cursor, err := repo.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return []model.Team{}, err
	}

	var teams []model.Team
	for cursor.Next(ctx) {
		var team model.Team
		err = cursor.Decode(&team)
		if err != nil {
			return teams, err
		}
		teams = append(teams, team)
	}
	return teams, nil
}

func (repo *mongoTeamRepository) AddAnEmployee(ctx context.Context, employeeId string, teamId string) error {
	_, err := repo.members.InsertOne(ctx, bson.D{
		{"employee_id", employeeId},
		{"team_id", teamId},
	})
	return err
}

func (repo *mongoTeamRepository) FindByTeamId(ctx context.Context, teamId string) ([]string, error) {
	cursor, err := repo.members.Find(ctx, bson.M{"team_id": teamId})
	if err != nil {
		return nil, err
	}

	var employeeList []string
	for cursor.Next(ctx) {
		var member Member
		err = cursor.Decode(&member)
		if err != nil {
			return employeeList, err
		}

		employeeList = append(employeeList, member.EmployeeId)
	}
	return employeeList, nil
}

func (repo *mongoTeamRepository) DeleteByTeamId(ctx context.Context, teamId string) error {
	_, err := repo.members.DeleteMany(ctx, bson.M{"team_id": teamId})
	return err
}

func (repo *mongoTeamRepository) DeleteAnEmployee(ctx context.Context, employeeId string, teamId string) error {
	_, err := repo.members.DeleteOne(ctx, bson.M{"employee_id": employeeId, "team_id": teamId})
	return err
}

func newMongoTeamRepository() (repo TeamRepository, err error) {
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

	repo = &mongoTeamRepository{
		client:     client,
		collection: client.Database(config.Get().Database).Collection(config.Get().Collection),
		members:    client.Database(config.Get().Database).Collection(config.Get().TeamMemberTable),
	}
	return repo, nil
}

func (repo *mongoTeamRepository) FindByUID(ctx context.Context, uid string) (model.Team, error) {
	result := repo.collection.FindOne(ctx, bson.M{"uid": uid})

	var team model.Team
	if err := result.Decode(&team); err != nil {
		return model.Team{}, err
	}
	return team, nil
}

func (repo *mongoTeamRepository) Save(ctx context.Context, team model.Team) error {
	_, err := repo.collection.InsertOne(ctx, team)
	return err
}

func (repo *mongoTeamRepository) Update(ctx context.Context, uid string, team model.Team) error {
	err := repo.collection.FindOneAndReplace(ctx, bson.M{"uid": uid}, team)
	return err.Err()
}

func (repo *mongoTeamRepository) Remove(ctx context.Context, uid string) error {
	_, err := repo.collection.DeleteOne(ctx, bson.M{"uid": uid})
	return err
}

type TeamDocument struct {
	ID          primitive.ObjectID `bson:"_id"`
	UID         string             `bson:"uid"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
}

func toTeamDocument(s model.Team) TeamDocument {
	return TeamDocument{
		ID:          primitive.NewObjectID(),
		UID:         s.UID,
		Name:        s.Name,
		Description: s.Description,
	}
}

func (s TeamDocument) ToModel() model.Team {
	return model.Team{
		UID:         s.UID,
		Name:        s.Name,
		Description: s.Description,
	}
}

type Member struct {
	ID         primitive.ObjectID `bson:"_id"`
	EmployeeId string             `bson:"employee_id"`
	TeamId     string             `bson:"team_id"`
}

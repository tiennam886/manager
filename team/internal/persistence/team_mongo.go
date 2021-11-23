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
	_, err := repo.collection.UpdateOne(ctx, bson.M{"uid": uid}, toTeamDocument(team))
	return err
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

package repositories

import (
	"context"
	"greens-basket/data"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repositorer interface {
	Insert(context.Context, *data.AppUser) error

	FindOne(ctx context.Context, filter primitive.M, result *data.AppUser) error

	GetAll(context.Context) ([]*data.AppUser, error)
	GetByID(context.Context, string) (*data.AppUser, error)

	ReplaceOne(ctx context.Context, id string, p *data.AppUser) error
	UpdateByID(ctx context.Context, subject string, doc primitive.M, result *data.AppUser) error
}

type UserRepository struct {
	DB         *mongo.Database
	Collection string
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		DB:         db,
		Collection: "users",
	}
}

func (s *UserRepository) Insert(ctx context.Context, p *data.AppUser) error {
	res, err := s.DB.Collection(s.Collection).InsertOne(ctx, p)
	if err != nil {
		return err
	}
	p.ID = res.InsertedID.(primitive.ObjectID).Hex()

	return err
}

func (s *UserRepository) GetAll(ctx context.Context) ([]*data.AppUser, error) {
	cursor, err := s.DB.Collection(s.Collection).Find(ctx, map[string]any{})
	if err != nil {
		return nil, err
	}

	products := []*data.AppUser{}
	err = cursor.All(ctx, &products)
	return products, err
}

func (s *UserRepository) GetByID(ctx context.Context, id string) (*data.AppUser, error) {
	var (
		objID, _ = primitive.ObjectIDFromHex(id)
		res      = s.DB.Collection(s.Collection).FindOne(ctx, bson.M{"_id": objID})
		p        = &data.AppUser{}
		err      = res.Decode(p)
	)
	return p, err
}

func (s *UserRepository) FindOne(ctx context.Context, filter primitive.M, result *data.AppUser) error {

	uCollection := s.DB.Collection(s.Collection)

	err := uCollection.FindOne(context.Background(), filter).Decode(result)
	if err != nil {
		return err
	}

	return nil

}

func (s *UserRepository) UpdateByID(ctx context.Context, subject string, doc primitive.M, result *data.AppUser) error {

	objectID, err := primitive.ObjectIDFromHex(subject)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID} // Replace "your_id_here" with the actual _id
	update := bson.M{"$set": doc}

	// Insert the BSON document into MongoDB
	return s.DB.Collection(s.Collection).FindOneAndUpdate(ctx, filter, update).Decode(result)

}

func (s *UserRepository) ReplaceOne(ctx context.Context, id string, p *data.AppUser) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Marshal the struct into a BSON map
	bsonDoc, err := bson.Marshal(p)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}

	// Insert the BSON document into MongoDB
	r, e := s.DB.Collection(s.Collection).ReplaceOne(ctx, filter, bsonDoc)

	if e != nil {
		return e
	}

	if r.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil

}

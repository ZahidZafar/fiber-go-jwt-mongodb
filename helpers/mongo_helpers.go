package helpers

import (
	"context"
	"greens-basket/data"
	"greens-basket/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoHelper struct {
	Client *mongo.Client
}

func (mh *MongoHelper) FindById(id string, result data.Collectioner) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	return mh.FindOne(filter, result)

}

func (mh *MongoHelper) ReplaceById(id string, result data.Collectioner) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	return mh.FindOne(filter, result)

}

func (mh *MongoHelper) FindOne(filter primitive.M, result data.Collectioner) error {

	db := mh.Client.Database(utils.Database)
	uCollection := db.Collection(result.Collection())

	err := uCollection.FindOne(context.Background(), filter).Decode(result)
	if err != nil {
		return err
	}

	return nil

}

func (mh *MongoHelper) InsertOne(val data.Collectioner) (string, error) {

	// Marshal the struct into a BSON map
	bsonDoc, err := bson.Marshal(val)
	if err != nil {
		return "", err
	}

	db := mh.Client.Database(utils.Database)
	uCollection := db.Collection(val.Collection())

	// Insert the BSON document into MongoDB
	r, e := uCollection.InsertOne(context.TODO(), bsonDoc)

	if e != nil {
		return "", e
	}

	return r.InsertedID.(primitive.ObjectID).Hex(), nil

}

func (mh *MongoHelper) UpdateByID(id string, ud primitive.M, collection string) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}

	update := bson.M{"$set": ud}

	db := mh.Client.Database(utils.Database)
	uCollection := db.Collection(collection)

	// Insert the BSON document into MongoDB
	_, e := uCollection.UpdateByID(context.TODO(), filter, update)

	if e != nil {
		return e
	}

	return nil

}

func (mh *MongoHelper) ReplaceOne(id string, val data.Collectioner) error {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}

	// Marshal the struct into a BSON map
	bsonDoc, err := bson.Marshal(val)
	if err != nil {
		return err
	}

	db := mh.Client.Database(utils.Database)
	uCollection := db.Collection(val.Collection())

	// Insert the BSON document into MongoDB
	_, e := uCollection.ReplaceOne(context.TODO(), filter, bsonDoc)

	if e != nil {
		return e
	}

	return nil

}

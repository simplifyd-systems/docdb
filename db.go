package docdb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DBIntf interface {
	Save(ctx context.Context, collection string, data interface{}) (string, error)
	SaveMultiple(context.Context, string, []interface{}) ([]interface{}, error)
	GetItem(ctx context.Context, collection string, filter map[string]interface{}, excludedFields map[string]interface{}, result interface{}) error
	GetItems(ctx context.Context, collection string, filter map[string]interface{}, limit int64, excludedFields map[string]interface{}, sort map[string]interface{}, results interface{}) error
	CountItems(ctx context.Context, collection string, filter map[string]interface{}) (int64, error)
	DeleteItem(ctx context.Context, c string, filter map[string]interface{}) (int64, error)
	DeleteItems(ctx context.Context, c string, filter map[string]interface{}) (int64, error)
	UpdateItem(ctx context.Context, c string, match map[string]interface{}, update map[string]interface{}) (int64, error)
	UpdateItems(ctx context.Context, c string, match map[string]interface{}, update map[string]interface{}) (int64, error)
	GetCollection(collection string) *mongo.Collection
	GetClient() *mongo.Client
}

// ErrMongoDBDuplicate error
var ErrMongoDBDuplicate = errors.New("duplicate entry")

// ErrInvalidObjectID error
var ErrInvalidObjectID = errors.New("invalid object ID")

// ErrNotFound error
var ErrNotFound = errors.New("item not found")

// MongoDB connection holder
type MongoDB struct {
	client   *mongo.Client
	database string
}

// NewDB creates a DB connection and returns a db instance
func NewDB(ctx context.Context, uri, database string) (db *MongoDB, err error) {
	db = &MongoDB{}

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return
	}
	err = client.Connect(ctx)
	if err != nil {
		return
	}

	db.database = database
	db.client = client
	return
}

// Disconnect closes the mongodb connection
func (db *MongoDB) Disconnect(ctx context.Context) {
	db.client.Disconnect(ctx)
}

// Ping db
func (db *MongoDB) Ping(ctx context.Context) (bool, error) {
	err := db.client.Ping(ctx, readpref.Primary())
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetClient func
func (db *MongoDB) GetClient() *mongo.Client {
	return db.client
}

// GetCollection func
func (db *MongoDB) GetCollection(collection string) *mongo.Collection {
	return db.client.Database(db.database).Collection(collection)
}

// Save func: c stands for collection where data would be saved. e.g save data in 'users' collection in MongoDB
// ctx can be a mongodb session context for transactions
func (db *MongoDB) Save(ctx context.Context, c string, data interface{}) (string, error) {
	collection := db.GetCollection(c)

	// ctx can be a mongodb session context for transactions
	insertResult, err := collection.InsertOne(ctx, data)
	if err != nil {
		/*
			var merr mongo.WriteException
			merr = err.(mongo.WriteException)
			errCode := merr.WriteErrors[0].Code
			if errCode == 11000 {
				return "", ErrMongoDBDuplicate
			} */
		return "", err
	}
	// update rule with returned ID
	return insertResult.InsertedID.(primitive.ObjectID).Hex(), nil
}

// SaveMultiple func: c stands for collection where data would be saved. e.g save data in 'users' collection in MongoDB
// ctx can be a mongodb session context for transactions
func (db *MongoDB) SaveMultiple(ctx context.Context, c string, items []interface{}) ([]interface{}, error) {
	collection := db.GetCollection(c)

	insertManyResult, err := collection.InsertMany(ctx, items)
	if err != nil {
		return nil, err
	}

	return insertManyResult.InsertedIDs, nil
}

// GetItem func: c stands for collection where item should be retrieved. e.g retrieve item from 'users' collection in MongoDB.
// ctx can be a mongodb session context for transactions
// results is a pointer to object to store returned data. nil is returned for error if item is found
func (db *MongoDB) GetItem(ctx context.Context, c string, filter map[string]interface{}, excludedFields map[string]interface{}, result interface{}) error {
	collection := db.GetCollection(c)

	findOptions := options.FindOne().SetProjection(excludedFields)

	// var result interface{}

	err := collection.FindOne(ctx, filter, findOptions).Decode(result)
	if err != nil {
		// TODO check for not found errror and return it
		// return nil, ErrNotFound
		return err
	}

	return nil
}

func (db *MongoDB) UpdateItem(ctx context.Context, c string, match map[string]interface{}, update map[string]interface{}) (int64, error) {
	collection := db.GetCollection(c)

	result, err := collection.UpdateOne(
		ctx,
		match,
		update,
	)
	if err != nil {
		return 0, err
	}

	return result.ModifiedCount, nil
}

func (db *MongoDB) UpdateItems(ctx context.Context, c string, match map[string]interface{}, update map[string]interface{}) (int64, error) {
	collection := db.GetCollection(c)

	result, err := collection.UpdateMany(
		ctx,
		match,
		update,
	)
	if err != nil {
		return 0, err
	}

	return result.ModifiedCount, nil
}

// GetItems func: c stands for collection where data would be saved. e.g save data in 'users' collection in MongoDB. id is string
// ctx can be a mongodb session context for transactions
// results is a pointer to slice of object to store returned data. nil is returned for error if item is found
func (db *MongoDB) GetItems(ctx context.Context, c string, filter map[string]interface{}, limit int64, excludedFields map[string]interface{}, sort map[string]interface{}, results interface{}) error {
	collection := db.GetCollection(c)

	findOptions := options.Find().SetProjection(excludedFields)
	findOptions.SetSort(sort)
	findOptions.SetLimit(limit)

	// var results []interface{}

	cur, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return err
	}
	// Close the cursor once finished
	defer cur.Close(ctx)

	if err := cur.All(ctx, results); err != nil {
		return err
	}

	return nil
}

// CountItems func: c stands for collection where items should be counted. e.g count items in 'users' collection in MongoDB.
// ctx can be a mongodb session context for transactions
func (db *MongoDB) CountItems(ctx context.Context, c string, filter map[string]interface{}) (int64, error) {
	collection := db.GetCollection(c)

	countOptions := options.Count()

	var result int64

	result, err := collection.CountDocuments(ctx, filter, countOptions)
	if err != nil {
		return 0, err
	}

	return result, nil
}

// DeleteItem func: c stands for collection where item should be retrieved. e.g retrieve item from 'users' collection in MongoDB.
// ctx can be a mongodb session context for transactions
func (db *MongoDB) DeleteItem(ctx context.Context, c string, filter map[string]interface{}) (int64, error) {
	collection := db.GetCollection(c)

	deleteResult, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}

	return deleteResult.DeletedCount, nil
}

// DeleteItems func: c stands for collection where item should be retrieved. e.g retrieve item from 'users' collection in MongoDB.
// ctx can be a mongodb session context for transactions
func (db *MongoDB) DeleteItems(ctx context.Context, c string, filter map[string]interface{}) (int64, error) {
	collection := db.GetCollection(c)

	deleteResult, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}

	return deleteResult.DeletedCount, nil
}

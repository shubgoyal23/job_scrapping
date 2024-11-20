package helpers

import (
	"context"
	"nScrapper/types"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDBConn *mongo.Client
var MongoDBName string

// inti mongodb
func InitMongoDB() error {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		LogError("connection mongodb", err)
		return err
	}
	MongoDBConn = client
	if err := MongoDBConn.Ping(context.Background(), nil); err != nil {
		LogError("connection mongodb", err)
		return err
	}
	MongoDBName = os.Getenv("MONGODB_NAME")
	return nil
}

// insert data in mongodb
func InsertMongoDB(data types.JobDataScrapeMap) error {
	if err := MongoDBConn.Ping(context.Background(), nil); err != nil {
		LogError("connection to mongodb failed with error: ", err)
		return err
	}
	collection := MongoDBConn.Database(MongoDBName).Collection("jobsScrapeMap")
	_, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		LogError("cannot insert in mongodb", err)
		return err
	}
	return nil
}

// get one data in mongodb
func GetOneDocMongoDB(filter bson.M) (interface{}, error) {
	var res interface{}
	if err := MongoDBConn.Ping(context.Background(), nil); err != nil {
		LogError("connection to mongodb failed with error: ", err)
		return res, err
	}
	collection := MongoDBConn.Database(MongoDBName).Collection("jobsScrapeMap")
	elem := collection.FindOne(context.Background(), filter)
	if elem.Err() != nil {
		LogError("cannot get doc in mongodb", elem.Err())
		return res, elem.Err()
	}
	if err := elem.Decode(&res); err != nil {
		LogError("cannot decode doc in mongodb", err)
		return res, err
	}
	return res, nil
}

// get one data in mongodb
func GetManyDocMongoDB(collectionName string, filter bson.M) ([]interface{}, error) {
	res := []interface{}{}
	if err := MongoDBConn.Ping(context.Background(), nil); err != nil {
		LogError("connection to mongodb failed with error: ", err)
		return res, err
	}
	collection := MongoDBConn.Database(MongoDBName).Collection(collectionName)
	elem, err := collection.Find(context.Background(), filter)
	if err != nil {
		LogError("cannot get doc in mongodb", err)
		return res, err
	}
	for elem.Next(context.Background()) {
		var result interface{}
		if err := elem.Decode(&result); err != nil {
			LogError("cannot decode doc in mongodb", err)
			return res, err
		}
		res = append(res, result)
	}
	return res, nil
}

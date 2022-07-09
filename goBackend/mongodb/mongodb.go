package mongodb

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongoClient *mongo.Client
var db *mongo.Database

func Init() {
	if client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", os.Getenv("MONGODBHOST"), os.Getenv("MONGODBPORT")))); err != nil {
		panic(err)
	} else {
		mongoClient = client
	}

	// Ping the primary
	if err := mongoClient.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	db = mongoClient.Database("msg")
	if err := CreateIndex("UserInfo", "username", true); err != nil {
		panic(err)
	}
	fmt.Println("db ready")
}

func Close() {
	if err := mongoClient.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

// CreateIndex - creates an index for a specific field in a collection
func CreateIndex(collectionName string, field string, unique bool) error {

	// 1. Lets define the keys for the index we want to create
	mod := mongo.IndexModel{
		Keys:    bson.M{field: 1}, // index in ascending order or -1 for descending order
		Options: options.Index().SetUnique(unique),
	}

	// 2. Create the context for this operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 3. Connect to the database and access the collection
	collection := db.Collection(collectionName)

	// 4. Create a single index
	_, err := collection.Indexes().CreateOne(ctx, mod)
	if err != nil {
		// 5. Something went wrong, we log it and return false
		fmt.Println(err.Error())
		return err
	}

	// 6. All went well, we return true
	return nil
}

type simpleCollection interface {
	Select() error
	Insert() (string, error)
	Update() error
	Delete() error
}

func getCollection(c simpleCollection) *mongo.Collection {
	return db.Collection(reflect.TypeOf(c).Elem().Name())
}

// func InsertOne(data interface{}) (string, error) {
// 	if result, err := db.Collection(reflect.TypeOf(data).Name()).InsertOne(context.TODO(), data); err != nil {
// 		return "", err
// 	} else {
// 		return result.InsertedID.(primitive.ObjectID).String(), nil
// 	}
// }

// func Find(data interface{}, result interface{}) (interface{}, error) {
// 	return result, db.Collection(reflect.TypeOf(data).Name()).FindOne(context.TODO(), data).Decode(&result)
// }

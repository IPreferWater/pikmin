package db

import (
	"context"
	"fmt"
	"os"

	"github.com/dktunited/api-output/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
//TODO give bomb as a lot of variant, unify the words
var (
	db *mongo.Client

	user         = os.Getenv("DB_USER")
	password     = os.Getenv("DB_PASSWORD")
	host         = os.Getenv("DB_HOST")
	port         = os.Getenv("DB_PORT")
	databaseName = os.Getenv("DB_NAME")
)

func InitDatabase() {
	//TODO
	//dbURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", user, password, host, port, databaseName)
	dbURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", "olimar", "password", "db-mongo", "27017", "pikmin-database")
	clientOptions := options.Client().ApplyURI(dbURI)

	var err error
	db, err = mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		panic(err)
	}

	err = db.Ping(context.TODO(), nil)

	if err != nil {
		panic(err)
	}
}

// CRUD
//create
func CreatePikmin(model model.Pikmin) (string, error) {
	collection := getCollection("pikmins")
	//generate a new ID
	model.ID = primitive.NewObjectID().String()
	_, err := collection.InsertOne(context.TODO(), model)
	return model.ID, err
}

//get
func GetPikminsByColor(color string) ([]model.Pikmin, error) {
	collection := getCollection("pikmins")
	filter := bson.M{"color": bson.M{"$eq": color}}
	res, err := collection.Find(context.TODO(), filter)

	var pikmins []model.Pikmin
	if err = res.All(context.TODO(), &pikmins); err != nil {
		return pikmins, err
	}
	return pikmins, nil
}

//update
func InsertPikmin(id string, newModel string) error {

	collection := getCollection("pikmins")

	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": bson.M{"$eq": id}}
	update := bson.M{"$set": newModel}

	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	return err
}

//upsert
//delete

//bulk update
func GiveBombs(pikmins []model.Pikmin) (int64, error) {
	collection := getCollection("pikmins")

	update := bson.M{"$set": bson.M{"weapons": bson.M{"bomb": true}}}
	var models []mongo.WriteModel
	for _, pikmin := range pikmins {
		filter := bson.M{"_id": bson.M{"$eq": pikmin.ID}}
		writeModel := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update)
		models = append(models, writeModel)
	}

	opts := options.BulkWrite().SetOrdered(false)
	res, err := collection.BulkWrite(context.TODO(), models, opts)
	return res.MatchedCount, err
}

func getCollection(tableName string) *mongo.Collection {
	//return db.Database(databaseName).Collection(tableName)
	return db.Database("pikmin-database").Collection(tableName)
}

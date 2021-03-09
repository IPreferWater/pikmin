package db

import (
	"context"
	"fmt"
	"os"

	"github.com/ipreferwater/pikmin/go/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
//TODO give bomb as a lot of variant, unify the words
var (
	user         = os.Getenv("DB_USER")
	password     = os.Getenv("DB_PASSWORD")
	host         = os.Getenv("DB_HOST")
	port         = os.Getenv("DB_PORT")
	databaseName = os.Getenv("DB_NAME")
)

type PikminRepoMongo struct {
	client *mongo.Client
}

func InitDatabase() {
	dbURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", user, password, host, port, databaseName)
	clientOptions := options.Client().ApplyURI(dbURI)

	db, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		panic(err)
	}

	err = db.Ping(context.TODO(), nil)

	if err != nil {
		panic(err)
	}

	PikminRepo = newMongoDDPikminRepository(db)
}

func newMongoDDPikminRepository(db *mongo.Client) PikminRepository {
	if db == nil {
		panic("missing db")
	}
	return &PikminRepoMongo{client: db}
}

// CRUD
//create
func (r *PikminRepoMongo) CreatePikmin(model model.Pikmin) (string, error) {
	collection := getCollection("pikmins", r.client)
	//generate a new ID
	model.ID = primitive.NewObjectID().String()
	_, err := collection.InsertOne(context.TODO(), model)
	return model.ID, err
}

//get
func (r *PikminRepoMongo) GetPikminsByColor(color string) ([]model.Pikmin, error) {
	collection := getCollection("pikmins", r.client)
	filter := bson.M{"color": bson.M{"$eq": color}}
	res, err := collection.Find(context.TODO(), filter)

	var pikmins []model.Pikmin
	if err = res.All(context.TODO(), &pikmins); err != nil {
		return pikmins, err
	}
	return pikmins, nil
}

//update
//TODO
func (r *PikminRepoMongo) UpdatePikmin(id string, newModel string) error {

	collection := getCollection("pikmins", r.client)

	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": bson.M{"$eq": id}}
	update := bson.M{"$set": newModel}

	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	return err
}

//upsert
//delete

//bulk update
func (r *PikminRepoMongo) GiveBombs(pikmins []model.Pikmin) (int64, error) {
	collection := getCollection("pikmins", r.client)

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

func getCollection(tableName string, client *mongo.Client) *mongo.Collection {
	return client.Database("pikmin-database").Collection(tableName)
}

package mongohelper

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	db *mongo.Database
}

func NewMongo(url string, collection string) *Mongo {
	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	fmt.Println("Connected to MongoDB!")
	db := client.Database(collection)
	return &Mongo{
		db: db,
	}
}

func (m *Mongo) Close() error {
	return m.db.Client().Disconnect(context.TODO())
}

func (m *Mongo) Insert(collection string, document []interface{}) error {
	result, err := m.db.Collection(collection).InsertMany(context.TODO(), document)
	if err != nil {
		return err
	}
	fmt.Println("Inserted : ", result.InsertedIDs)
	return nil
}

func (m *Mongo) GetCollection(collection string) (*mongo.Cursor, error) {
	c, err := m.db.Collection(collection).Find(context.TODO(), bson.D{}, options.Find())
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (m *Mongo) GetDocument(collection string, document interface{}) (*mongo.Cursor, error) {
	r, err := m.db.Collection(collection).Find(context.TODO(), document)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (m *Mongo) UpdateDocument(collection string, filter interface{}, update interface{}) error {
	result, err := m.db.Collection(collection).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	fmt.Println("updated document", result.MatchedCount)
	return nil
}

func (m *Mongo) DeleteDocument(collection string, filter interface{}) error {
	result, err := m.db.Collection(collection).DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	fmt.Println("deleted document", result.DeletedCount)
	return nil
}

package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	logs "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Connection struct {
	Conn     *mongo.Client
	ConnDB   *mongo.Database
	Database string
}

// Client - MongoDB Connection
var Client *Connection

// NewConnection - new connection of amqp
func NewConnection() error {

	

	mongoUrl := "mongodb://127.0.0.1:27017"
	mongoDatabase := "pr1"

	if mongoUrl == "" || mongoDatabase == "" {
		logs.Error("configuration is missing for mongodb")
		return errors.New("configuration is missing for mongodb")
	}

	mongoClient := &Connection{
		Conn:     nil,
		ConnDB:   nil,
		Database: mongoDatabase,
	}
	ctx := context.Background()
	clientOptions := options.Client().ApplyURI(mongoUrl)


	var err error
	mongoClient.Conn, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		logs.Error("cannot connect new client to db")
		return err
	}
	err = mongoClient.Conn.Ping(context.TODO(), nil)
	if err != nil {
		logs.Error(err)
		return err
	}

	fmt.Println("Connected to MongoDB2")

	mongoClient.ConnDB = mongoClient.Conn.Database(mongoClient.Database)
	Client = mongoClient

	return nil
}

// GetCollection - Helper Functions
func GetCollection(collectionName string) *mongo.Collection {
	return Client.ConnDB.Collection(collectionName)
}

// DbContext - Helper Functions
func DbContext(i time.Duration) (context.Context, context.CancelFunc) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), i * time.Second)
	return ctx, ctxCancel
}



package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"time"

	"service2/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbHost = "localhost"
	port   = ":1234"
)

type App struct {
	client *mongo.Client
}

func (a *App) CreateUser(user *models.User, reply *string) error {
	coll := a.client.Database("taskDB").Collection("users")
	doc := bson.D{{"email", user.Email}, {"salt", user.Salt}, {"password", user.Password}}
	_, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		return err
	}
	*reply = ""
	return nil
}

func (a *App) GetUser(email string, user *models.User) error {
	coll := a.client.Database("taskDB").Collection("users")
	doc := bson.D{{"email", email}}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	if err := coll.FindOne(ctx, doc).Decode(&user); err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	return nil
}

func main() {
	dbURI := fmt.Sprintf("mongodb://%s:27017", dbHost)
	client, err := mongo.NewClient(options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	app := new(App)

	app.client = client
	rpc.Register(app)
	t, err := net.ResolveTCPAddr("tcp4", port)
	if err != nil {
		fmt.Println(err)
		return
	}
	l, err := net.ListenTCP("tcp4", t)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		c, err := l.Accept()
		if err != nil {
			continue
		}
		rpc.ServeConn(c)
	}
}

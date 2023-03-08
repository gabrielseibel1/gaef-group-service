package main

import (
	"context"
	"flag"
	"fmt"
	"gaef-group-service/handler"
	"gaef-group-service/service"
	"gaef-group-service/store"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type handlers struct {
	handler Handler
}

type Handler interface {
	Handle() gin.HandlerFunc
}

func main() {
	// read command-line args
	var prod bool
	flag.BoolVar(&prod, "production", false, "indicates the service is used for production")
	flag.Parse()

	// read environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	dbURI := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("MONGODB_DATABASE")
	collectionName := os.Getenv("MONGODB_COLLECTION")

	// TODO: secure connection to mongo with user/password
	// connect to mongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}

	// instantiate and inject dependencies
	store := store.New(client.Database(dbName).Collection(collectionName))
	service := service.New(store)
	handler := handler.New(service)
	handlers := handlers{
		handler: handler,
	}

	// run http server
	r := gin.Default()
	users := r.Group("/api/v0/groups")
	{
		public := users.Group("")
		{
			public.POST("/", handlers.handler.Handle())
		}
		auth := users.Group("", handlers.handler.Handle())
		{
			auth.GET("/:id", handlers.handler.Handle())
		}
	}
	r.Run(fmt.Sprintf("0.0.0.0:%s", port))
}

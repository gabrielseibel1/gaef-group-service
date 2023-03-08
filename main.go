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
		auth:          handler,
		onlyLeaders:   handler,
		createGroup:   handler,
		readAllGroups: handler,
		readGroup:     handler,
		updateGroup:   handler,
		deleteGroup:   handler,
		readMembers:   handler,
		addMember:     handler,
		deleteMember:  handler,
		readLeaders:   handler,
		addLeader:     handler,
		deleteLeader:  handler,
	}

	// run http server
	server := gin.Default()
	api := server.Group("/api/v0/groups", handlers.auth.AuthMiddleware())
	{
		api.POST("/", handlers.createGroup.CreateGroupHandler())
		api.GET("/", handlers.readAllGroups.ReadAllGroupsHandler())
		api.GET("/:id", handlers.readGroup.ReadGroupHandler())
		api.GET("/:id/members", handlers.readMembers.ReadMembersHandler())
		api.GET("/:id/leaders", handlers.readLeaders.ReadLeadersHandler())

		forLeaders := api.Group("", handlers.onlyLeaders.OnlyLeadersMiddleware())
		{
			forLeaders.PUT("/:id", handlers.updateGroup.UpdateGroupHandler())
			forLeaders.DELETE("/:id", handlers.deleteGroup.DeleteGroupHandler())
			forLeaders.POST("/:id/members", handlers.addMember.AddMemberHandler())
			forLeaders.DELETE(":gid/members/:mid", handlers.deleteMember.DeleteMemberHandler())
			forLeaders.POST("/:id/leaders", handlers.addLeader.AddLeadersHandler())
			forLeaders.DELETE(":gid/leaders/:lid", handlers.deleteLeader.DeleteLeaderHandler())
		}
	}
	server.Run(fmt.Sprintf("0.0.0.0:%s", port))
}

type handlers struct {
	auth          AuthMiddleware
	onlyLeaders   OnlyLeadersMiddleware
	createGroup   CreateGroupHandler
	readAllGroups ReadAllGroupsHandler
	readGroup     ReadGroupHandler
	updateGroup   UpdateGroupHandler
	deleteGroup   DeleteGroupHandler
	readMembers   ReadMembersHandler
	addMember     AddMemberHandler
	deleteMember  DeleteMemberHandler
	readLeaders   ReadLeadersHandler
	addLeader     AddLeadersHandler
	deleteLeader  DeleteLeaderHandler
}

type AuthMiddleware interface {
	AuthMiddleware() gin.HandlerFunc
}
type OnlyLeadersMiddleware interface {
	OnlyLeadersMiddleware() gin.HandlerFunc
}
type CreateGroupHandler interface {
	CreateGroupHandler() gin.HandlerFunc
}
type ReadAllGroupsHandler interface {
	ReadAllGroupsHandler() gin.HandlerFunc
}
type ReadGroupHandler interface {
	ReadGroupHandler() gin.HandlerFunc
}
type UpdateGroupHandler interface {
	UpdateGroupHandler() gin.HandlerFunc
}
type DeleteGroupHandler interface {
	DeleteGroupHandler() gin.HandlerFunc
}
type ReadMembersHandler interface {
	ReadMembersHandler() gin.HandlerFunc
}
type AddMemberHandler interface {
	AddMemberHandler() gin.HandlerFunc
}
type DeleteMemberHandler interface {
	DeleteMemberHandler() gin.HandlerFunc
}
type ReadLeadersHandler interface {
	ReadLeadersHandler() gin.HandlerFunc
}
type AddLeadersHandler interface {
	AddLeadersHandler() gin.HandlerFunc
}
type DeleteLeaderHandler interface {
	DeleteLeaderHandler() gin.HandlerFunc
}

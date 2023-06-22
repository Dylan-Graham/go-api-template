package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DB_NAME = "flyFantasyLeague"
const ATHLETE_COLLECTION = "athletes"
const USER_COLLECTION = "user"

func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func main() {
	client, ctx := connectDB()

	router := gin.Default()

	Init_Routes(router, client, ctx)

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func connectDB() (*mongo.Client, context.Context) {
	uri := goDotEnvVariable("MONGO_URI")

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		log.Fatal("Had a database connection error: ", err)
	}

	log.Default().Println("Successfuly connected to database")

	return client, ctx
}

func Init_Routes(router *gin.Engine, client *mongo.Client, ctx context.Context) {
	router.GET("/user", func(c *gin.Context) {
		users := fetchUsers(client, ctx)

		c.JSON(http.StatusOK, gin.H{
			"users": users,
		})
	})
}

func fetchUsers(client *mongo.Client, ctx context.Context) []string {
	userCollection := client.Database(DB_NAME).Collection(USER_COLLECTION)
	filter := bson.D{}
	cur, err := userCollection.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)
	users := []string{}
	for cur.Next(ctx) {
		var result bson.D
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		// do something with result....
		for _, elem := range result {
			if elem.Key == "name" {
				nameValue := elem.Value.(string)
				users = append(users, nameValue)
				break
			}
		}
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return users
}

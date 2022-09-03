package config

import(
	"context"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Ctx = context.Background()

func Connect() (*mongo.Database, error) {
	clientOptions := options.Client()
	godotenv.Load(".env")
	mongoUrl := os.Getenv("MONGO_URL")
	clientOptions.ApplyURI(mongoUrl)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
			return nil, err
	}

	err = client.Connect(Ctx)
	if err != nil {
			return nil, err
	}

	return client.Database("go-chat-app"), nil
}
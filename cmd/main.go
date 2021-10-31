package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/microphoneabuser/auth-service/pkg/handler"
	"github.com/microphoneabuser/auth-service/pkg/repository"
	"github.com/microphoneabuser/auth-service/pkg/service"
	"github.com/microphoneabuser/auth-service/server"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initialization config: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	mongoClient, err := repository.NewMongoClient(
		repository.MongoConfig{
			Uri:      viper.GetString("mongo.Uri"),
			Username: viper.GetString("mongo.Username"),
			Password: os.Getenv("MONGO_PASSWORD"),
		})
	if err != nil {
		log.Fatalf("error initialization mongo: %s", err.Error())
	}

	db := mongoClient.Database(viper.GetString("mongo.DBName"))
	log.Println("connection to mongo is ready")

	repos := repository.NewRepository(db)
	service := service.NewService(repos)
	handlers := handler.NewHandler(service)

	if err := server.RunServer(viper.GetString("port"), handlers); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

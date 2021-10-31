package service

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/microphoneabuser/auth-service/models"
	"github.com/microphoneabuser/auth-service/pkg/repository"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Authorization interface {
	GenerateTokens(ctx context.Context, id primitive.ObjectID) (Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
	CreateUser(ctx context.Context) (models.User, error)
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Repository) *Service {
	accessTokenTTL, refreshTokenTTL := getTTLs()
	return &Service{
		Authorization: NewAuthService(repos.Authorization, time.Duration(accessTokenTTL)*time.Hour, time.Duration(refreshTokenTTL)*time.Hour),
	}
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

func getTTLs() (int, int) {
	accessTokenTTL, err := strconv.Atoi(viper.GetString("jwt.accessTokenTTL"))
	if err != nil {
		log.Fatalf("Error reading accessTokenTTL from env var: %s", err)
	}
	refreshTokenTTL, err := strconv.Atoi(viper.GetString("jwt.refreshTokenTTL"))
	if err != nil {
		log.Fatalf("Error reading refreshTokenTTL from env var: %s", err)
	}
	return accessTokenTTL, refreshTokenTTL
}

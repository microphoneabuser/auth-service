package repository

import (
	"context"

	"github.com/microphoneabuser/auth-service/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Authorization interface {
	GetById(ctx context.Context, id primitive.ObjectID) (models.User, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (models.User, error)
	SetSession(ctx context.Context, id primitive.ObjectID, session models.Session) error
	CreateUser(ctx context.Context) (models.User, error)
}

type Repository struct {
	Authorization
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		Authorization: NewAuthMongo(db),
	}
}

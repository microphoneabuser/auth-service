package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/microphoneabuser/auth-service/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthMongo struct {
	db *mongo.Collection
}

func NewAuthMongo(db *mongo.Database) *AuthMongo {
	return &AuthMongo{db: db.Collection(authCollection)}
}

func (r *AuthMongo) SetSession(ctx context.Context, id primitive.ObjectID, session models.Session) error {
	_, err := r.db.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"session": session, "lastVisitAt": time.Now()}})

	return err
}

func (r *AuthMongo) GetByRefreshToken(ctx context.Context, refreshToken string) (models.User, error) {
	var user models.User
	if err := r.db.FindOne(ctx, bson.M{
		"session.refreshToken": refreshToken,
		"session.expiresAt":    bson.M{"$gt": time.Now()},
	}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.User{}, models.ErrorUserNotFound
		}

		return models.User{}, err
	}

	return user, nil
}

func (r *AuthMongo) GetById(ctx context.Context, id primitive.ObjectID) (models.User, error) {
	var user models.User
	if err := r.db.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.User{}, models.ErrorUserNotFound
		}

		return models.User{}, err
	}

	return user, nil
}

func (r *AuthMongo) CreateUser(ctx context.Context) (models.User, error) {
	user := models.User{
		ID: primitive.NewObjectID(),
	}
	res, err := r.db.InsertOne(ctx, &user)
	if err != nil {
		if isDuplicate(err) {
			return models.User{}, models.ErrorUserAlreadyExists
		}
		return models.User{}, err
	}

	log.Println(res.InsertedID)

	user.ID = res.InsertedID.(primitive.ObjectID)

	return user, nil
}

func isDuplicate(err error) bool {
	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == 11000 {
				return true
			}
		}
	}

	return false
}

package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/microphoneabuser/auth-service/models"
	"github.com/microphoneabuser/auth-service/pkg/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo            repository.Authorization
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewAuthService(repo repository.Authorization, accessTokenTTL time.Duration, refreshTokenTTL time.Duration) *AuthService {
	return &AuthService{
		repo:            repo,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (s *AuthService) GenerateTokens(ctx context.Context, id primitive.ObjectID) (Tokens, error) {
	user, err := s.repo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, models.ErrorUserNotFound) {
			return Tokens{}, err
		}
		return Tokens{}, err
	}

	return s.createSession(ctx, user.ID)
}

func (s *AuthService) RefreshTokens(ctx context.Context, id primitive.ObjectID, refreshToken string) (Tokens, error) {
	user, err := s.repo.GetById(ctx, id)
	if err != nil {
		return Tokens{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Session.RefreshToken), []byte(refreshToken))
	if err != nil {
		return Tokens{}, models.ErrorWrongRefreshToken
	}

	return s.createSession(ctx, user.ID)
}

func (s *AuthService) CreateUser(ctx context.Context) (models.User, error) {
	return s.repo.CreateUser(ctx)
}

func (s *AuthService) createSession(ctx context.Context, id primitive.ObjectID) (Tokens, error) {
	var (
		res Tokens
		err error
	)

	res.AccessToken, err = newJWT(id.Hex(), s.accessTokenTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = newRefreshToken()
	if err != nil {
		return res, err
	}
	refreshTokenBcrypted, err := bcryptToken(res.RefreshToken)
	if err != nil {
		return res, err
	}

	session := models.Session{
		RefreshToken: refreshTokenBcrypted,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}

	err = s.repo.SetSession(ctx, id, session)

	return res, err
}

func newJWT(id string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(ttl).Unix(),
		Subject:   id,
	})

	return token.SignedString([]byte(os.Getenv("SIGNING_KEY")))
}

func newRefreshToken() (string, error) {
	bytes := make([]byte, 32)

	source := rand.NewSource(time.Now().Unix())
	random := rand.New(source)

	if _, err := random.Read(bytes); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", bytes), nil
}

func bcryptToken(token string) (string, error) {
	var bytes []byte
	var err error
	if bytes, err = bcrypt.GenerateFromPassword([]byte(token), 10); err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (s *AuthService) ParseToken(accessToken string) (string, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SIGNING_KEY")), nil
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", claims["sub"]), nil
}

package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

const DB_NAME string = "users-db"
const COLLECTION_NAME string = "users"

type Service struct {
	collection *mongo.Collection
}

func newService(client *mongo.Client) Service {
	return Service{
		collection: client.Database(DB_NAME).Collection(COLLECTION_NAME),
	}
}

func validateCredentials(creds Credentials) error {
	switch {
	case len(creds.Username) == 0:
		return errors.New("username must not be empty")
	case strings.Contains(creds.Username, " "):
		return errors.New("username must not contain whitespace")
	case len(creds.Password) < 8:
		return errors.New("password must have at least 8 characters")
	case strings.Contains(creds.Password, " "):
		return errors.New("password must not contain whitespace")
	default:
		return nil
	}
}

func (s Service) findUser(username string) (*User, error) {
	var user User
	err := s.collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	return &user, err
}

func (s Service) userExists(username string) (bool, error) {
	_, err := s.findUser(username)
	switch {
	case err == nil:
		return true, nil
	case err == mongo.ErrNoDocuments:
		return false, nil
	default:
		return false, err
	}
}

func (s Service) createUser(creds Credentials) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = s.collection.InsertOne(context.TODO(), NewUser(
		primitive.NewObjectID(),
		creds.Username,
		string(hashedPassword),
	))
	return err
}

func passwordMatches(inputPassword string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
}

func createToken(username string, signingKey []byte) (string, error) {
	expirationTime := time.Now().Add(24 * 7 * time.Hour)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey)
}

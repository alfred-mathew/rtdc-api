package auth

import (
	"context"
	"errors"
	"strings"

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
	signingKey []byte
}

func newService(client *mongo.Client, signingKey []byte) Service {
	return Service{
		collection: client.Database(DB_NAME).Collection(COLLECTION_NAME),
		signingKey: signingKey,
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

func (s Service) createToken(username string) (string, error) {
	claims := &Claims{
		Username: username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.signingKey)
}

func (s Service) parseClaimsFromAuthHeader(header string) (*Claims, error) {
	if header == "" {
		return nil, errors.New("empty authorization header")
	}

	if !strings.HasPrefix(header, "Bearer") {
		return nil, errors.New("bearer token not present in authorization header")
	}

	splits := strings.Split(header, " ")
	if len(splits) != 2 {
		return nil, errors.New("malformed bearer token")
	}

	tokenString := splits[1]
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return s.signingKey, nil
	})
	return claims, err
}

package mongo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rizzmatch/rm_auth/internal/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	db *mongo.Client
}

func New() (*Storage, error) {
	const op = "storage.mongo.New"

	opts := options.Client().ApplyURI("mongodb+srv://revoltik:revoltik@cluster0.zxojcem.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0")

	db, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(context.TODO(), nil); err != nil {
		panic(err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	return s.db.Disconnect(context.TODO())
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (uid uuid.UUID, err error) {
	const op = "storage.mongo.SaveUser"

	usersCollection := s.db.Database("rizzmatch").Collection("users")

	uuidv6, err := uuid.NewV6()
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	// TODO: add validation
	id, err := usersCollection.InsertOne(ctx, models.User{
		ID:       uuidv6,
		Email:    email,
		PassHash: passHash,
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	objectId := id.InsertedID.(primitive.Binary).Data

	lastInserdUUID, err := uuid.FromBytes(objectId)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return lastInserdUUID, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.mongo.User"

	var tmpUser = new(models.User)

	usersCollection := s.db.Database("rizzmatch").Collection("users")
	err := usersCollection.FindOne(ctx, bson.M{"email": email}).Decode(tmpUser)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return *tmpUser, nil
}

func (s *Storage) IsAdmin(ctx context.Context, email string) (bool, error) {
	const op = "storage.mongo.isAdmin"

	tmpUser := new(models.User)
	usersCollection := s.db.Database("rizzmatch").Collection("users")
	err := usersCollection.FindOne(ctx, bson.M{"email": email}).Decode(tmpUser)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if email == "admin" {
		return true, nil
	}
	return false, nil
}

package mongo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rizzmatch/rm_auth/internal/core/models"
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

	fmt.Println(id)
	fmt.Println(id.InsertedID)

	return id.InsertedID.(uuid.UUID), nil
}

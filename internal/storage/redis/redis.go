package redis

// import (
// 	"context"

// 	"github.com/google/uuid"
// 	"github.com/redis/go-redis/v9"
// 	"github.com/rizzmatch/rm_auth/internal/core/models"
// )

// type Storage struct {
// 	db *redis.Client
// }

// func New() (*Storage, error) {
// 	const op = "storage.redis.New"

// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:     "localhost:6379",
// 		Password: "",
// 		DB:       0,
// 	})

// 	_ = op

// 	err := rdb.Ping(context.TODO()).Err()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &Storage{db: rdb}, nil
// }

// func (s *Storage) Close() error {
// 	return s.db.Close()
// }

// func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (uuid.UUID, error) {
// 	const op = "storage.redis.SaveUser"
// 	_ = op

// 	uid, err := uuid.NewV6()
// 	if err != nil {
// 		return uuid.Nil, err
// 	}

// 	err = s.db.Set(ctx, email, passHash, 0).Err()
// 	if err != nil {
// 		return uuid.Nil, err
// 	}

// 	return uid, nil
// }

// func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
// 	const op = "storage.redis.User"
// 	_ = op

// 	res, err := s.db.Get(ctx, email).Result()
// 	if err != nil {
// 		return models.User{}, nil
// 	}

// 	uid, err := uuid.NewV6()
// 	if err != nil {
// 		return models.User{}, nil
// 	}

// 	tmpUser := models.User{
// 		ID:       uid,
// 		Email:    email,
// 		PassHash: []byte(res),
// 	}

// 	return tmpUser, nil
// }

// func (s *Storage) IsAdmin(ctx context.Context, email string) (bool, error) {
// 	return true, nil
// }

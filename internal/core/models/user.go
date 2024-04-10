package models

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `bson:"_id"`
	Email    string    `bson:"email"`
	PassHash []byte    `bson:"pass_hash"`
}

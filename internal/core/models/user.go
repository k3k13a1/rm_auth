package models

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `bson:"_id" json:"id"`
	Email    string    `bson:"email" json:"email"`
	PassHash []byte    `bson:"pass_hash" json:"pass_hash"`
}

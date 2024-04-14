package models

type User struct {
	ID       int    `bson:"_id" json:"id"`
	Email    string `bson:"email" json:"email"`
	PassHash []byte `bson:"pass_hash" json:"pass_hash"`
}

package models

type UserModel struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
	Level uint `bson:"level,omitempty"`
}
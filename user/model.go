package user

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id       primitive.ObjectID `bson:"_id" json:"id"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"-"`
	Name     string             `bson:"name" json:"name"`
	Age      int                `bson:"age" json:"age"`
}

type Users []User

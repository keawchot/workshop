package user

import (
	"context"
	"time"
	"workshop/db"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	resource   *db.Resource
	collection *mongo.Collection
}

type Repository interface {
	GetAll() (Users, error)
	CreateOne(UserRequest) (User, error)
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByEmailAndPassword(email string, password string) (*User, error)
	UpdateUser(id string, userUpdateRequest UserUpdateRequest) (*User, error)
}

// NewUserRepository create repository
func NewUserRepository(resource *db.Resource) Repository {
	collection := resource.DB.Collection("user")
	repository := &UserRepository{resource: resource, collection: collection}
	return repository
}

// GetAll to get all users
func (ur *UserRepository) GetAll() (Users, error) {
	users := Users{}
	ctx, cancel := initContext()
	defer cancel()

	cursor, err := ur.collection.Find(ctx, bson.M{})
	if err != nil {
		return Users{}, err
	}

	for cursor.Next(ctx) {
		var user User
		err = cursor.Decode(&user)
		if err != nil {
			logrus.Print(err)
		}
		users = append(users, user)
	}
	return users, nil
}

// GetByID to get user by id
func (ur *UserRepository) GetByID(id string) (*User, error) {
	var user User

	ctx, cancel := initContext()
	defer cancel()

	objID, _ := primitive.ObjectIDFromHex(id)
	err := ur.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByID to get user by email
func (ur *UserRepository) GetByEmail(email string) (*User, error) {
	var user User

	ctx, cancel := initContext()
	defer cancel()

	err := ur.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) GetByEmailAndPassword(email string, password string) (*User, error) {
	var user User

	ctx, cancel := initContext()
	defer cancel()

	err := ur.collection.FindOne(ctx, bson.M{"email": email, "password": password}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// CreateOne to create new user
func (ur *UserRepository) CreateOne(userRequest UserRequest) (User, error) {
	user := User{
		Id:       primitive.NewObjectID(),
		Email:    userRequest.Email,
		Password: userRequest.Password,
		Name:     userRequest.Name,
		Age:      userRequest.Age,
	}
	ctx, cancel := initContext()
	defer cancel()
	_, err := ur.collection.InsertOne(ctx, user)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (ur *UserRepository) UpdateUser(id string, userUpdateRequest UserUpdateRequest) (*User, error) {
	var user User
	ctx, cancel := initContext()
	defer cancel()
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{primitive.E{Key: "_id", Value: objID}}

	update := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "password", Value: userUpdateRequest.Password},
		primitive.E{Key: "name", Value: userUpdateRequest.Name},
		primitive.E{Key: "age", Value: userUpdateRequest.Age},
	}}}

	err := ur.collection.FindOneAndUpdate(ctx, filter, update).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func initContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	return ctx, cancel
}

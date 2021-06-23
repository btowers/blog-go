package mongodb

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/btowers/blog-go/pkg/adder"
	"github.com/btowers/blog-go/pkg/auth"
	"github.com/btowers/blog-go/pkg/lister"
	"github.com/btowers/blog-go/pkg/remover"
	"github.com/btowers/blog-go/pkg/updater"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ExistingUser indicates that the email is already registered
	ErrEmailAlreadyRegistered = errors.New("email address is already registered")

	// ErrReadingUserFromDB indicates that there was a problem while getting the user from the DB
	ErrReadingUserFromDB = errors.New("error reading user from db")

	// ErrEmailPassword indicates that the email or password dont match
	ErrEmailPassword = errors.New("wrong email/password")

	// ErrSavingUserInDB indicates that the there was an error while inserting a User in the DB
	ErrSavingUserInDB = errors.New("wrong email/password")

	// ErrSavingPostInDB indicates that the there was an error while inserting a Post in the DB
	ErrSavingPostInDB = errors.New("error saving Post in DB")
)

type Storage struct {
	db *mongo.Database
}

func NewStorage() *Storage {
	s := new(Storage)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	s.db = client.Database("aurube")
	return (s)
}

// Authentication

func (s *Storage) Register(u auth.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if user is already registered
	var userFound bson.M = nil
	filter := bson.M{"email": u.Email}
	s.db.Collection("users").FindOne(ctx, filter).Decode(&userFound)
	if userFound != nil {
		return ErrEmailAlreadyRegistered
	}

	// Password encryption
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	u.Password = string(hash)

	// Save user in DB with encrypted password
	userInserted, err := s.db.Collection("users").InsertOne(ctx, u)
	if err != nil {
		return ErrSavingUserInDB
	}
	fmt.Println(userInserted)
	return nil
}

func (s *Storage) Login(u auth.User) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var userFound interface{}
	filter := bson.M{"email": u.Email}
	err := s.db.Collection("users").FindOne(ctx, filter).Decode(&userFound)
	if err == mongo.ErrNoDocuments {
		return nil, mongo.ErrNoDocuments
	} else if err != nil {
		return nil, ErrReadingUserFromDB
	}

	bsonUser, _ := bson.Marshal(userFound)

	var usr auth.User
	bson.Unmarshal(bsonUser, &usr)

	var busr bson.M
	bson.Unmarshal(bsonUser, &busr)

	errs := bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(u.Password))
	if errs != nil {
		return nil, ErrEmailPassword
	}
	return busr, nil
}

func (s *Storage) Logout() error {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return nil
}

func (s *Storage) IsAuthenticated() error {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return nil
}

func (s *Storage) GetUser(email string) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var userFound interface{}
	filter := bson.M{"email": email}
	err := s.db.Collection("users").FindOne(ctx, filter).Decode(&userFound)
	if err == mongo.ErrNoDocuments {
		return nil, mongo.ErrNoDocuments
	} else if err != nil {
		return nil, ErrReadingUserFromDB
	}

	bsonUser, _ := bson.Marshal(userFound)

	var busr bson.M
	bson.Unmarshal(bsonUser, &busr)

	return busr, nil
}

func (s *Storage) DeleteUser(u remover.User) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := s.db.Collection("users")
	collection.FindOneAndDelete(ctx, bson.M{"email": u.Email})
}

func (s *Storage) UpdateUser(email string, u updater.User) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := s.db.Collection("users")
	filter := bson.M{"email": email}
	update := bson.M{"$set": u}
	opts := options.FindOneAndUpdate().SetUpsert(false)
	var updatedDocument bson.M
	err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedDocument)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return
		}
		log.Fatal(err)
	}
	fmt.Printf("UPDATED USER ID: %v \n", updatedDocument["_id"])
}

// CRUD Post

func (s *Storage) AddPost(p adder.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Save user in DB with encrypted password
	postInserted, err := s.db.Collection("posts").InsertOne(ctx, p)
	if err != nil {
		return ErrSavingPostInDB
	}
	fmt.Println(postInserted)
	return nil
}

func (s *Storage) GetPost(p lister.Post) (lister.Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	postsCollection := s.db.Collection("posts")
	objectId, err := primitive.ObjectIDFromHex(p.Id)
	if err != nil {
		return lister.Post{}, mongo.ErrInvalidIndexValue
	}

	filter := bson.M{"_id": objectId}
	var postFound lister.Post
	postsCollection.FindOne(ctx, filter).Decode(&postFound)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return lister.Post{}, mongo.ErrNoDocuments
		}
	}
	return postFound, err
}

func (s *Storage) DeletePost(p remover.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := s.db.Collection("posts")
	collection.FindOneAndDelete(ctx, bson.M{"Id": p.Id})
	return nil
}

func (s *Storage) UpdatePost(p updater.Post) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection := s.db.Collection("posts")
	filter := bson.M{"Id": p.Id}
	update := bson.M{"$set": p}
	opts := options.FindOneAndUpdate().SetUpsert(false)
	var updatedDocument bson.M
	err := collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedDocument)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return mongo.ErrNoDocuments
		}
		log.Fatal(err)
	}
	fmt.Printf("UPDATED USER ID: %v \n", updatedDocument["_id"])
	return nil
}

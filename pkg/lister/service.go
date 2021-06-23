package lister

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	GetUser(string) (interface{}, error)
	GetPost(Post) (Post, error)
}

type Service interface {
	GetUser(string) (interface{}, error)
	GetPost(Post) (Post, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s service) GetUser(email string) (interface{}, error) {
	s.r.GetUser(email)
	return nil, nil
}

func (s service) GetPost(p Post) (Post, error) {
	post, err := s.r.GetPost(p)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return Post{}, mongo.ErrNoDocuments
		}
	}
	return post, nil
}

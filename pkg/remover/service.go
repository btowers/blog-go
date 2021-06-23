package remover

type Repository interface {
	DeleteUser(User)
	DeletePost(Post) error
}

type Service interface {
	DeleteUser(User)
	DeletePost(Post) error
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s service) DeleteUser(u User) {
	s.r.DeleteUser(u)

}

func (s service) DeletePost(p Post) error {
	s.r.DeletePost(p)
	return nil
}

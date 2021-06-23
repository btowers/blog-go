package updater

type Repository interface {
	UpdateUser(string, User)
	UpdatePost(Post) error
}

type Service interface {
	UpdateUser(string, User)
	UpdatePost(Post) error
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s service) UpdateUser(email string, u User) {
	s.r.UpdateUser(email, u)

}

func (s service) UpdatePost(p Post) error {
	s.r.UpdatePost(p)
	return nil
}

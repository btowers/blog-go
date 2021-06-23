package adder

type Repository interface {
	AddPost(Post) error
}

type Service interface {
	AddPost(Post) error
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s service) AddPost(p Post) error {
	s.r.AddPost(p)
	return nil
}

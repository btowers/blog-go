package auth

type Repository interface {
	Register(User) error
	Login(User) (interface{}, error)
	Logout() error
	IsAuthenticated() error
	GetUser(string) (interface{}, error)
}

type Service interface {
	Register(User) error
	Login(User) (interface{}, error)
	Logout() error
	IsAuthenticated() error
	GetUser(string) (interface{}, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s service) Register(u User) error {
	s.r.Register(u)
	return nil
}

func (s service) Login(u User) (interface{}, error) {
	user, err := s.r.Login(u)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s service) Logout() error {
	s.r.Logout()
	return nil
}

func (s service) IsAuthenticated() error {
	s.r.IsAuthenticated()
	return nil
}

func (s service) GetUser(email string) (interface{}, error) {
	user, err := s.r.GetUser(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

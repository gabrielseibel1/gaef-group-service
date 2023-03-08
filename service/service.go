package service

type Service struct {
	storer Storer
}

func New(storer Storer) Service {
	return Service{
		storer: storer,
	}
}

func (s Service) Serve() {
	s.storer.Store()
}

type Storer interface {
	Store()
}

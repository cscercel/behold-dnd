package character

import (
	"context"
)


type Service interface {
	ListCharacters(ctx context.Context) (error)
}

type svc struct {

}

func NewService() Service {
	return &svc{}
}


func (s *svc) ListCharacters(ctx context.Context) error {
	return nil
}

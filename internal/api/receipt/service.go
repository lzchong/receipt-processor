package receipt

import (
	"errors"
)

type Service interface {
	Points(id string) (int64, error)
	SetPoints() string
}

type serviceImpl struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &serviceImpl{repository}
}

var ErrReceiptNotFound = errors.New("receipt not found")

func (s *serviceImpl) Points(id string) (int64, error) {
	points, ok := s.repository.Points(id)
	if !ok {
		return 0, ErrReceiptNotFound
	}
	return points, nil
}

func (s *serviceImpl) SetPoints() string {
	return "7fb1377b-b223-49d9-a31a-5a02701dd310"
}

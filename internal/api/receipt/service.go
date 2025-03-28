package receipt

import (
	"errors"
)

type Service interface {
	Points(id string) (int64, error)
	Process(receipt *Receipt) string
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

func (s *serviceImpl) Process(receipt *Receipt) string {
	points := receipt.CalculatePoints()
	id := s.repository.CreatePoints(points)
	return id
}

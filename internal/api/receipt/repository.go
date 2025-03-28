package receipt

import (
	"sync"
)

type Repository interface {
	Points(id string) (int64, bool)
}

type inMemoryRepository struct {
	lock   sync.RWMutex
	points map[string]int64
}

func NewRepository() Repository {
	return &inMemoryRepository{
		points: make(map[string]int64),
	}
}

func (s *inMemoryRepository) Points(id string) (int64, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	points, ok := s.points[id]
	return points, ok
}

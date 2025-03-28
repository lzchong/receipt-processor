package receipt

import (
	"sync"

	"github.com/google/uuid"
)

type Repository interface {
	Points(id string) (int64, bool)
	CreatePoints(points int64) string
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

func (s *inMemoryRepository) CreatePoints(points int64) string {
	s.lock.Lock()
	defer s.lock.Unlock()

	id := s.generateID()
	s.points[id] = points
	return id
}

func (s *inMemoryRepository) generateID() string {
	for {
		id := uuid.New().String()
		_, ok := s.points[id]
		if !ok {
			return id
		}
	}
}

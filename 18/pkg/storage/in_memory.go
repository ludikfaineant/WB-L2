package storage

import (
	"fmt"
	"sync"
	"time"
	"wb_l12/18/internal/model"
)

type Storage interface {
	Create(event *model.Event) (int, error)
	Update(event *model.Event) error
	Delete(id int) error
	GetByDay(user_id int, date time.Time) ([]model.Event, error)
	GetByWeek(user_id int, date time.Time) ([]model.Event, error)
	GetByMonth(user_id int, date time.Time) ([]model.Event, error)
}

var ErrNotFound = fmt.Errorf("event not found")

type InMemoryStorage struct {
	events map[int]model.Event
	mu     sync.RWMutex
	nextID int
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		events: make(map[int]model.Event),
		nextID: 1,
	}
}

func (s *InMemoryStorage) Create(event *model.Event) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.nextID
	event.ID = id
	s.events[id] = *event
	s.nextID++
	return id, nil
}

func (s *InMemoryStorage) Update(event *model.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[event.ID]; !ok {
		return ErrNotFound
	}
	s.events[event.ID] = *event

	return nil
}

func (s *InMemoryStorage) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[id]; !ok {
		return ErrNotFound
	}
	delete(s.events, id)

	return nil
}

func (s *InMemoryStorage) GetByDay(user_id int, date time.Time) ([]model.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var res []model.Event
	for _, e := range s.events {
		if e.UserID == user_id && isSameDay(e.Date, date) {
			res = append(res, e)
		}
	}
	return res, nil
}

func (s *InMemoryStorage) GetByWeek(user_id int, date time.Time) ([]model.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var res []model.Event
	for _, e := range s.events {
		if e.UserID == user_id && isSameWeek(e.Date, date) {
			res = append(res, e)
		}
	}
	return res, nil

}
func (s *InMemoryStorage) GetByMonth(user_id int, date time.Time) ([]model.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var res []model.Event
	for _, e := range s.events {
		if e.UserID == user_id && isSameMonth(e.Date, date) {
			res = append(res, e)
		}
	}
	return res, nil

}

func isSameDay(a, b time.Time) bool {
	return a.Day() == b.Day() && a.Month() == b.Month() && a.Year() == b.Year()
}

func isSameWeek(a, b time.Time) bool {
	y1, w1 := a.ISOWeek()
	y2, w2 := b.ISOWeek()
	return y1 == y2 && w1 == w2
}

func isSameMonth(a, b time.Time) bool {
	return a.Month() == b.Month() && a.Year() == b.Year()
}

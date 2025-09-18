package service

import (
	"time"
	"wb_l12/18/internal/model"
	"wb_l12/18/pkg/storage"
)

type Service struct {
	storage storage.Storage
}

func NewService(storage storage.Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) CreateEvent(userID int, date time.Time, title string) (int, error) {
	event := &model.Event{
		UserID: userID,
		Date:   date,
		Title:  title,
	}
	return s.storage.Create(event)
}

func (s *Service) UpdateEvent(id, userID int, date time.Time, title string) error {
	event := &model.Event{
		ID:     id,
		UserID: userID,
		Date:   date,
		Title:  title,
	}
	return s.storage.Update(event)
}

func (s *Service) DeleteEvent(id int) error {
	return s.storage.Delete(id)
}

func (s *Service) GetByDay(userID int, date time.Time) ([]model.Event, error) {
	return s.storage.GetByDay(userID, date)
}

func (s *Service) GetByWeek(userID int, date time.Time) ([]model.Event, error) {
	return s.storage.GetByWeek(userID, date)
}
func (s *Service) GetByMonth(userID int, date time.Time) ([]model.Event, error) {
	return s.storage.GetByMonth(userID, date)
}

package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"wb_l12/18/pkg/storage"
)

func TestCreateEvent_Success(t *testing.T) {
	service := NewService(storage.NewInMemoryStorage())
	id, err := service.CreateEvent(1, time.Now(), "Test")
	assert.NoError(t, err)
	assert.Greater(t, id, 0)
}

func TestGetByDay_ReturnsMatchingEvents(t *testing.T) {
	service := NewService(storage.NewInMemoryStorage())
	now := time.Now()

	service.CreateEvent(1, now, "Event 1")
	service.CreateEvent(1, now, "Event 2")
	service.CreateEvent(2, now, "Other user")

	events, _ := service.GetByDay(1, now)

	assert.Len(t, events, 2)
}

func TestGetByDay_WithNoEvents(t *testing.T) {
	service := NewService(storage.NewInMemoryStorage())
	events, err := service.GetByDay(1, time.Now())
	assert.NoError(t, err)
	assert.Empty(t, events)
}

func TestUpdateEvent_Success(t *testing.T) {
	storage := storage.NewInMemoryStorage()
	service := NewService(storage)

	id, _ := service.CreateEvent(1, time.Now(), "Old Title")

	err := service.UpdateEvent(id, 1, time.Now(), "New Title")
	assert.NoError(t, err)

	events, _ := service.GetByDay(1, time.Now())
	assert.Equal(t, "New Title", events[0].Title)
}

func TestDeleteEvent_RemovesEvent(t *testing.T) {
	storage := storage.NewInMemoryStorage()
	service := NewService(storage)

	id, _ := service.CreateEvent(1, time.Now(), "To delete")

	err := service.DeleteEvent(id)
	assert.NoError(t, err)

	events, _ := service.GetByDay(1, time.Now())
	assert.Empty(t, events)
}

func TestGetByWeek_IncludesEventsFromSameWeek(t *testing.T) {
	storage := storage.NewInMemoryStorage()
	service := NewService(storage)

	wednesday := time.Date(2023, 12, 27, 0, 0, 0, 0, time.UTC)
	tuesday := time.Date(2023, 12, 26, 0, 0, 0, 0, time.UTC)

	service.CreateEvent(1, tuesday, "Meeting")
	service.CreateEvent(1, wednesday, "Party")

	events, err := service.GetByWeek(1, wednesday)
	assert.NoError(t, err)
	assert.Len(t, events, 2)
}

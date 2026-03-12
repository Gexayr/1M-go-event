package events

import (
	"encoding/json"
	"errors"
	"strconv"

	"gorm.io/gorm"
)

type Service interface {
	GetEventDetails(id uint) (*EventDetailsResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetEventDetails(id uint) (*EventDetailsResponse, error) {
	event, err := s.repo.GetEventByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	clientID, _ := strconv.ParseUint(event.ClientID, 10, 64)

	return &EventDetailsResponse{
		ID:        event.ID,
		ClientID:  uint(clientID),
		EventType: event.EventType,
		RiskScore: event.RiskScore,
		Metadata:  json.RawMessage(event.Metadata),
		Timestamp: event.Timestamp,
	}, nil
}

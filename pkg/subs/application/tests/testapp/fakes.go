package testapp

import (
	"context"
	"sync"
)

// SpyEventPublisher сохраняет опубликованные события для тестов
type SpyEventPublisher struct {
	mu     sync.Mutex
	Events []PublishedEvent
}

// PublishedEvent хранит данные об одном событии
type PublishedEvent struct {
	Topic   string
	Payload []byte
}

// Реализация интерфейса application.EventPublisher
func (s *SpyEventPublisher) Publish(ctx context.Context, topic string, payload []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Events = append(s.Events, PublishedEvent{
		Topic:   topic,
		Payload: payload,
	})
	return nil
}

// Получить копию всех опубликованных событий (для безопасного чтения)
func (s *SpyEventPublisher) GetEvents() []PublishedEvent {
	s.mu.Lock()
	defer s.mu.Unlock()
	copied := make([]PublishedEvent, len(s.Events))
	copy(copied, s.Events)
	return copied
}

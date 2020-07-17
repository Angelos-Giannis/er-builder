package service

import "github.com/Angelos-Giannis/erd-builder/internal/domain"

// Service describes the service flow.
type Service struct {
	options domain.Options
}

// New creates and returns a new service.
func New(options domain.Options) *Service {
	return &Service{
		options: options,
	}
}

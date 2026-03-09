package service

import (
	"github.com/iamsainadh/neovault/audit/internal/repository"
	"github.com/twmb/franz-go/pkg/kgo"
)

type AuditService interface {
	
}

type auditService struct {
	repo   repository.AuditRepository
	client *kgo.Client
}

func NewConsumer(client *kgo.Client, repo repository.AuditRepository) AuditService {
	return &auditService{
		client: client,
		repo:   repo,
	}
}

package hubsvc

import (
	"context"

	"go-server/internal/config"
	"go-server/internal/models"
	"go-server/internal/repositories/contracts"
)

type Store = contracts.HubStore
type Operational interface {
	EnsureSocietyOperational(context.Context, int64) error
}
type Service struct {
	repo        Store
	operational Operational
	storage     models.ImageStorage
	prefix      string
	limits      config.HubConfig
	slots       chan struct{}
}

func New(repo Store, operational Operational, storage models.ImageStorage, prefix string, limits config.HubConfig) *Service {
	if repo == nil {
		panic("hubSvc: required dependency is nil")
	}
	return &Service{repo: repo, operational: operational, storage: storage, prefix: prefix, limits: limits, slots: make(chan struct{}, models.ImageConcurrency)}
}

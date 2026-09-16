package service

import (
	"searchengine/internal/config"
	"searchengine/internal/store"
	"searchengine/pkg/logger"
)

// Service 业务服务层。
type Service struct {
	store store.Store
	log   *logger.Logger
	cfg   *config.Config
}

// New 构造业务服务。
func New(st store.Store, log *logger.Logger, cfg *config.Config) *Service {
	return &Service{store: st, log: log, cfg: cfg}
}

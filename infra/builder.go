package infra

import (
	"context"
	"service/build"
	"service/config"
)

type builder struct {
	ctx context.Context
	cfg *config.Config
}

// NewBuilder новый сборщик инфраструктурных объектов.
func NewBuilder(ctx context.Context, cfg *config.Config) *builder {
	cfg.Service.Info.Version = build.Version

	return &builder{
		ctx: ctx,
		cfg: cfg,
	}
}

func (b *builder) GetInfo() *config.ServiceInfo {
	return b.cfg.Service.Info
}

package product

import (
	"context"
	"github.com/online-shop/internal/entity"
	"github.com/online-shop/pkg/log"
)

type Service interface {
	Get(ctx context.Context, id string) (entity.Product, error)
}

type service struct {
	repo   Repository
	logger log.Logger
}

// NewService creates a new album service.
func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

func (s service) Get(ctx context.Context, id string) (entity.Product, error) {
	product, err := s.repo.Get(ctx, id)
	if err != nil {
		return entity.Product{}, err
	}
	return product, nil
}

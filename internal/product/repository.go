package product

import (
	"context"
	"fmt"
	"github.com/online-shop/internal/entity"
	"github.com/online-shop/pkg/log"
	"github.com/online-shop/pkg/mysql"
)

type Repository interface {
	Get(ctx context.Context, id string) (entity.Product, error)
	List(ctx context.Context) ([]entity.Product, error)
}

// repository persists albums in database
type repository struct {
	db     mysql.BaseRepository
	logger log.Logger
}

// NewRepository creates a new album repository
func NewRepository(db mysql.BaseRepository, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) List(ctx context.Context) ([]entity.Product, error) {
	q := fmt.Sprintf("select id, name, stock, price from product")

	var products []entity.Product

	err := r.db.FetchRows(ctx, q, &products)
	if err != nil {
		return products, err
	}

	return products, nil
}

func (r repository) Get(ctx context.Context, id string) (entity.Product, error) {
	q := fmt.Sprintf("select id, name, stock, price from product where id = ?")

	var product entity.Product

	err := r.db.FetchRow(ctx, q, &product, id)
	if err != nil {
		return product, err
	}

	return product, nil
}

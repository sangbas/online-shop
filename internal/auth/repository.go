package auth

import (
	"context"
	"fmt"
	"github.com/online-shop/internal/entity"
	"github.com/online-shop/pkg/log"
	"github.com/online-shop/pkg/mysql"
)

type Repository interface {
	CreateUser(ctx context.Context, user entity.User) error
	FindByUsername(ctx context.Context, username string) (entity.User, error)
}

// repository persists users in database
type repository struct {
	db     mysql.BaseRepository
	logger log.Logger
}

// NewRepository creates a new users repository
func NewRepository(db mysql.BaseRepository, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) CreateUser(ctx context.Context, user entity.User) error {
	q := fmt.Sprintf("insert into user (id, username, fullname, phone, email, password, token) " +
		"values (:id, :username, :fullname, :phone, :email, :password, :token)")

	_, err := r.db.Exec(ctx, q, user)
	if err != nil {
		return err
	}

	return nil
}

func (r repository) FindByUsername(ctx context.Context, username string) (entity.User, error) {
	q := fmt.Sprintf("select * from user where username = ?")

	var user entity.User

	err := r.db.FetchRow(ctx, q, &user, username)
	if err != nil {
		return user, err
	}

	return user, nil
}

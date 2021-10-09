package mysql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

// BaseRepository type
type BaseRepository struct {
	MasterDB *sqlx.DB
	SlaveDB  *sqlx.DB
}

func (r *BaseRepository) Exec(ctx context.Context, query string, args interface{}) (sql.Result, error) {
	var (
		res sql.Result
		err error
	)

	if r.MasterDB == nil {
		return res, errors.New("the master DB connection is nil")
	}

	res, err = r.MasterDB.NamedExecContext(ctx, query, args)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// FetchRow the fetch data row on Slave DB
func (r *BaseRepository) FetchRows(query string, resp interface{}, args ...interface{}) error {
	if r.SlaveDB == nil {
		return errors.New("the slave DB connection is nil")
	}

	err := r.SlaveDB.Select(resp, query, args...)
	if err != nil {
		return err
	}

	return nil
}

// FetchRow the fetch data row on Slave DB
func (r *BaseRepository) FetchRow(query string, resp interface{}, args ...interface{}) error {
	if r.SlaveDB == nil {
		return errors.New("the slave DB connection is nil")
	}

	err := r.SlaveDB.Get(resp, query, args...)
	if err != nil {
		return err
	}

	return nil
}

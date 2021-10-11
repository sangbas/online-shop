package mysql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"go.uber.org/multierr"
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
func (r *BaseRepository) FetchRows(ctx context.Context, query string, resp interface{}, args ...interface{}) error {
	if r.SlaveDB == nil {
		return errors.New("the slave DB connection is nil")
	}

	err := r.SlaveDB.SelectContext(ctx, resp, query, args...)
	if err != nil {
		return err
	}

	return nil
}

// FetchRow the fetch data row on Slave DB
func (r *BaseRepository) FetchRow(ctx context.Context, query string, resp interface{}, args ...interface{}) error {
	if r.SlaveDB == nil {
		return errors.New("the slave DB connection is nil")
	}

	err := r.SlaveDB.GetContext(ctx, resp, query, args...)
	if err != nil {
		return err
	}

	return nil
}

//type TxFn func(Transaction) error
//
//func (r *BaseRepository) WithTransaction(db *sql.DB, fn TxFn) (err error) {
//	tx, err := db.Begin()
//	if err != nil {
//		return
//	}
//
//	defer func() {
//		if p := recover(); p != nil {
//			// a panic occurred, rollback and repanic
//			tx.Rollback()
//			panic(p)
//		} else if err != nil {
//			// something went wrong, rollback
//			tx.Rollback()
//		} else {
//			// all good, commit
//			err = tx.Commit()
//		}
//	}()
//
//	err = fn(tx)
//	return err
//}

func (r *BaseRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	tx, err := r.MasterDB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return tx, err
	}

	return tx, nil
}

func (r *BaseRepository) EndTx(tx *sql.Tx, err error) error {
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return multierr.Combine(err, rollbackErr)
		}

		return err
	} else {
		if commitErr := tx.Commit(); commitErr != nil {
			return err
		}

		return nil
	}
}

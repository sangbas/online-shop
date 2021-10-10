package order

import (
	"context"
	"fmt"
	"github.com/online-shop/internal/entity"
	"github.com/online-shop/pkg/log"
	"github.com/online-shop/pkg/mysql"
	"time"
)

type Repository interface {
	Get(ctx context.Context, id string) (entity.Order, error)
	PlaceOrder(ctx context.Context, orderReq entity.Order) error
	CreateOrder(ctx context.Context, order entity.Order) error
	CreateOrderDetail(ctx context.Context, orderDetail entity.OrderDetail) error
	UpdateOrder(ctx context.Context, order entity.Order) error
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

func (r repository) Get(ctx context.Context, id string) (entity.Order, error) {
	q := fmt.Sprintf("select * from orders where id = ?")

	var order entity.Order

	err := r.db.FetchRow(ctx, q, &order, id)
	if err != nil {
		return order, err
	}

	return order, nil
}

func (r repository) PlaceOrder(ctx context.Context, orderReq entity.Order) error {
	// Begin Transaction
	tx, err := r.db.BeginTx(ctx)
	defer func() {
		r.db.EndTx(tx, err)
	}()

	now := time.Now()

	err = r.CreateOrder(ctx, entity.Order{
		ID:        orderReq.ID,
		UserID:    orderReq.UserID,
		AddressID: orderReq.AddressID,
		OrderDate: &now,
		Status:    orderReq.Status,
		Amount:    orderReq.Amount,
	})

	for _, orderDetail := range orderReq.OrderDetails {
		err = r.CreateOrderDetail(ctx, entity.OrderDetail{
			ID:        entity.GenerateID(),
			OrderID:   orderReq.ID,
			ProductID: orderDetail.ProductID,
			Price:     orderDetail.Price,
			Quantity:  orderDetail.Quantity,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (r repository) CreateOrder(ctx context.Context, order entity.Order) error {
	q := fmt.Sprintf("insert into orders (id, user_id, address_id, order_date, status, amount) " +
		"values (:id, :user_id, :address_id, :order_date, :status, :amount)")

	_, err := r.db.Exec(ctx, q, order)
	if err != nil {
		return err
	}

	return nil
}

func (r repository) CreateOrderDetail(ctx context.Context, orderDetail entity.OrderDetail) error {
	q := fmt.Sprintf("insert into order_detail values (:id, :order_id, :product_id, :quantity, :price)")

	_, err := r.db.Exec(ctx, q, orderDetail)
	if err != nil {
		return err
	}

	return nil
}

func (r repository) UpdateOrder(ctx context.Context, order entity.Order) error {
	q := fmt.Sprintf("update orders set address_id = :address_id, " +
		"payment_date = :payment_date, " +
		"delivered_date = :delivered_date, " +
		"status = :status " +
		"where id = :id")

	_, err := r.db.Exec(ctx, q, order)
	if err != nil {
		return err
	}

	return nil
}

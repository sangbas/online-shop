package order

import (
	"context"
	"github.com/online-shop/internal/auth"
	"github.com/online-shop/internal/entity"
	"github.com/online-shop/pkg/log"
	"time"
)

const (
	WAITING_PAYMENT = "WAITING_FOR_PAYMENT"
	SEND_PAYMENT    = "SEND_PAYMENT"
	PAID            = "PAID"
	DELIVERED       = "DELIVERED"
	RECEIVED        = "RECEIVED"
	CANCELLED       = "CANCELLED"
	REJECTED        = "REJECTED"
)

type Service interface {
	Get(ctx context.Context, id string) (entity.Order, error)
	PlaceOrder(ctx context.Context, input PlaceOrderRequest) (entity.Order, error)
	UpdateOrder(ctx context.Context, input UpdateOrderRequest) (entity.Order, error)
}

type PlaceOrderRequest struct {
	ShippingAddress string        `json:"shipping_address"`
	Items           []ItemRequest `json:"items"`
}

type ItemRequest struct {
	ProductID int64   `json:"product_id"`
	Quantity  int32   `json:"quantity"`
	Price     float64 `json:"price"`
}

type UpdateOrderRequest struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

type service struct {
	repo   Repository
	logger log.Logger
}

// NewService creates a new album service.
func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

func (s service) Get(ctx context.Context, id string) (entity.Order, error) {
	order, err := s.repo.Get(ctx, id)
	if err != nil {
		return entity.Order{}, err
	}
	return order, nil
}

func (s service) PlaceOrder(ctx context.Context, input PlaceOrderRequest) (entity.Order, error) {
	var orderDetails []entity.OrderDetail
	var total float64

	orderId := entity.GenerateID()
	user := auth.CurrentUser(ctx)

	for _, item := range input.Items {
		orderDetails = append(orderDetails, entity.OrderDetail{
			ProductID: item.ProductID,
			Price:     item.Price,
			Quantity:  item.Quantity,
		})
		total += total + (item.Price * float64(item.Quantity))
	}

	err := s.repo.PlaceOrder(ctx, entity.Order{
		ID:           orderId,
		UserID:       user.GetID(),
		AddressID:    input.ShippingAddress,
		Status:       WAITING_PAYMENT,
		Amount:       total,
		OrderDetails: orderDetails,
	})

	if err != nil {
		return entity.Order{}, err
	}

	return s.Get(ctx, orderId)
}

func (s service) UpdateOrder(ctx context.Context, input UpdateOrderRequest) (entity.Order, error) {
	order, err := s.repo.Get(ctx, input.OrderID)
	if err != nil {
		return entity.Order{}, err
	}

	now := time.Now()
	switch status := input.Status; status {
	case SEND_PAYMENT:
		order.PaymentDate = &now
	case DELIVERED:
		order.DeliveredDate = &now
	}
	order.Status = input.Status

	err = s.repo.UpdateOrder(ctx, order)
	if err != nil {
		return entity.Order{}, err
	}

	return order, nil
}

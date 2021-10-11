package order

import (
	"context"
	"github.com/online-shop/internal/auth"
	"github.com/online-shop/internal/entity"
	"github.com/online-shop/pkg/log"
	"time"
)

const (
	CREATED   = "CREATED"
	PAYMENT   = "PAYMENT"
	VERIFIED  = "VERIFIED"
	SHIPPED   = "SHIPPED"
	RECEIVED  = "RECEIVED"
	CANCELLED = "CANCELLED"
	REJECTED  = "REJECTED"
)

type Service interface {
	Get(ctx context.Context, id string) (OrderResponse, error)
	PlaceOrder(ctx context.Context, input PlaceOrderRequest) (OrderResponse, error)
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

type ItemResponse struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int32   `json:"quantity"`
}

type OrderResponse struct {
	ID     string         `json:"id"`
	UserID string         `json:"user_id"`
	Status string         `json:"status"`
	Amount float64        `json:"amount"`
	Items  []ItemResponse `json:"items"`
}

type service struct {
	repo   Repository
	logger log.Logger
}

// NewService creates a new album service.
func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

func (s service) Get(ctx context.Context, id string) (OrderResponse, error) {
	order, err := s.repo.GetCompleteOrder(ctx, id)
	if err != nil {
		return OrderResponse{}, err
	}

	var items []ItemResponse
	for _, item := range order {
		items = append(items, ItemResponse{
			Name:     item.ProductName,
			Price:    item.Price,
			Quantity: item.Quantity,
		})
	}

	return OrderResponse{
		ID:     order[0].ID,
		UserID: order[0].UserID,
		Status: order[0].Status,
		Amount: order[0].Amount,
		Items:  items,
	}, nil
}

func (s service) PlaceOrder(ctx context.Context, input PlaceOrderRequest) (OrderResponse, error) {
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
		total = total + (item.Price * float64(item.Quantity))
	}

	err := s.repo.PlaceOrder(ctx, entity.Order{
		ID:           orderId,
		UserID:       user.GetID(),
		AddressID:    input.ShippingAddress,
		Status:       CREATED,
		Amount:       total,
		OrderDetails: orderDetails,
	})

	if err != nil {
		return OrderResponse{}, err
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
	case PAYMENT:
		order.PaymentDate = &now
	case SHIPPED:
		order.DeliveredDate = &now
	}
	order.Status = input.Status

	err = s.repo.UpdateOrder(ctx, order)
	if err != nil {
		return entity.Order{}, err
	}

	return order, nil
}

package order

import (
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/online-shop/internal/errors"
	"github.com/online-shop/pkg/log"
	"net/http"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler, logger log.Logger) {
	res := resource{service, logger}
	r.Use(authHandler)

	r.Post("/order", res.placeOrder)
	r.Put("/order", res.updateOrder)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) placeOrder(c *routing.Context) error {
	var input PlaceOrderRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}
	order, err := r.service.PlaceOrder(c.Request.Context(), input)
	if err != nil {
		return err
	}

	return c.WriteWithStatus(order, http.StatusCreated)
}

func (r resource) updateOrder(c *routing.Context) error {
	var input UpdateOrderRequest

	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}
	payment, err := r.service.UpdateOrder(c.Request.Context(), input)
	if err != nil {
		return err
	}

	return c.WriteWithStatus(payment, http.StatusOK)
}
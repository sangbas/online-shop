package order

import (
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/online-shop/internal/errors"
	"github.com/online-shop/internal/response"
	"github.com/online-shop/pkg/log"
	"net/http"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler, logger log.Logger) {
	res := resource{service, logger}
	r.Use(authHandler)

	r.Get("/orders/<id>", res.getOrder)
	r.Post("/orders", res.placeOrder)
	r.Put("/orders", res.updateOrder)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) getOrder(c *routing.Context) error {
	order, err := r.service.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.Write(order)
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
	_, err := r.service.UpdateOrder(c.Request.Context(), input)
	if err != nil {
		return err
	}

	return c.WriteWithStatus(response.SuccessResponse(), http.StatusOK)
}

package product

import (
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/online-shop/pkg/log"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler, logger log.Logger) {
	res := resource{service, logger}
	r.Use(authHandler)

	r.Get("/products", res.list)
	r.Get("/products/<id>", res.get)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) get(c *routing.Context) error {
	product, err := r.service.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.Write(product)
}

func (r resource) list(c *routing.Context) error {
	product, err := r.service.List(c.Request.Context())
	if err != nil {
		return err
	}

	return c.Write(product)
}

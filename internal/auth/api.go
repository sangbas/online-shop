package auth

import (
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/online-shop/internal/errors"
	"github.com/online-shop/internal/response"
	"github.com/online-shop/pkg/log"
	"net/http"
)

// RegisterHandlers registers handlers for different HTTP requests.
func RegisterHandlers(rg *routing.RouteGroup, service Service, logger log.Logger) {
	rg.Post("/login", login(service, logger))
	rg.Post("/register", register(service, logger))
}

// login returns a handler that handles user login request.
func login(service Service, logger log.Logger) routing.Handler {
	return func(c *routing.Context) error {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.Read(&req); err != nil {
			logger.With(c.Request.Context()).Errorf("invalid request: %v", err)
			return errors.BadRequest("")
		}

		token, err := service.Login(c.Request.Context(), req.Username, req.Password)
		if err != nil {
			return err
		}
		return c.Write(struct {
			Token string `json:"token"`
		}{token})
	}
}

// login returns a handler that handles user login request.
func register(service Service, logger log.Logger) routing.Handler {
	return func(c *routing.Context) error {
		var input RegisterRequest
		if err := c.Read(&input); err != nil {
			logger.With(c.Request.Context()).Info(err)
			return errors.BadRequest("")
		}
		err := service.CreateUser(c.Request.Context(), input)
		if err != nil {
			return err
		}

		return c.WriteWithStatus(response.SuccessResponse(), http.StatusCreated)
	}

}

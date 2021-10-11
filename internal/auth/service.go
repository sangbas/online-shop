package auth

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/online-shop/internal/entity"
	"github.com/online-shop/internal/errors"
	"github.com/online-shop/pkg/log"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// Service encapsulates the authentication logic.
type Service interface {
	// authenticate authenticates a user using username and password.
	// It returns a JWT token if authentication succeeds. Otherwise, an error is returned.
	Login(ctx context.Context, username, password string) (string, error)
	CreateUser(ctx context.Context, user RegisterRequest) error
}

// Identity represents an authenticated user identity.
type Identity interface {
	// GetID returns the user ID.
	GetID() string
	// GetName returns the user name.
	GetUsername() string
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"fullname"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

type service struct {
	signingKey      string
	tokenExpiration int
	logger          log.Logger
	repo            Repository
}

// NewService creates a new authentication service.
func NewService(repo Repository, signingKey string, tokenExpiration int, logger log.Logger) Service {
	return service{signingKey, tokenExpiration, logger, repo}
}

// Login authenticates a user and generates a JWT token if authentication succeeds.
// Otherwise, an error is returned.
func (s service) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	if identity := s.authenticate(ctx, user, password); identity != nil {
		return s.generateJWT(identity)
	}
	return "", errors.Unauthorized("")
}

func (s service) CreateUser(ctx context.Context, user RegisterRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = s.repo.CreateUser(ctx, entity.User{
		ID:       entity.GenerateID(),
		Username: user.Username,
		FullName: user.FullName,
		Phone:    user.Phone,
		Email:    user.Email,
		Password: string(hashedPassword),
		Token:    "",
	})

	if err != nil {
		return err
	}

	return nil
}

func verifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// authenticate authenticates a user using username and password.
// If username and password are correct, an identity is returned. Otherwise, nil is returned.
func (s service) authenticate(ctx context.Context, user entity.User, password string) Identity {
	logger := s.logger.With(ctx, "user", user.Username)

	err := verifyPassword(password, user.Password)
	if err != nil {
		logger.Infof("authentication failed")
		return nil
	}
	logger.Infof("authentication successful")
	return entity.User{ID: user.ID, Username: user.Username}
}

// generateJWT generates a JWT that encodes an identity.
func (s service) generateJWT(identity Identity) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   identity.GetID(),
		"name": identity.GetUsername(),
		"exp":  time.Now().Add(time.Duration(s.tokenExpiration) * time.Hour).Unix(),
	}).SignedString([]byte(s.signingKey))
}

package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/PolyAbit/auth/internal/lib/jwt"
	"github.com/PolyAbit/auth/internal/lib/logger/sl"
	"github.com/PolyAbit/auth/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type UserStorage interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error)
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Auth struct {
	log         *slog.Logger
	userStorage UserStorage
	tokenTTL    time.Duration
	tokenSecret string
}

func New(log *slog.Logger, userStorage UserStorage, tokenTTL time.Duration, tokenSecret string) *Auth {
	return &Auth{
		log:         log,
		userStorage: userStorage,
		tokenTTL:    tokenTTL,
	}
}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "services.auth.IsAdmin"

	log := a.log.With(slog.String("op", op), slog.Int64("userId", userID))

	isAdmin, err := a.userStorage.IsAdmin(ctx, userID)

	if errors.Is(err, models.ErrUserNotFound) {
		return false, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	if err != nil {
		log.Error("failed to get permission", err)

		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}

func (a *Auth) Login(ctx context.Context, email string, password string) (string, error) {
	const op = "services.auth.Login"

	log := a.log.With(slog.String("op", op), slog.String("email", email))

	user, err := a.userStorage.User(ctx, email)

	if errors.Is(err, models.ErrUserNotFound) {
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	if err != nil {
		log.Error("failed to get user", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	token, err := jwt.NewToken(user, a.tokenSecret, a.tokenTTL)

	if err != nil {
		log.Error("failed generate toke", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *Auth) RegisterNewUser(ctx context.Context, email string, password string) (userID int64, err error) {
	const op = "services.auth.RegisterNewUser"

	log := a.log.With(slog.String("op", op), slog.String("email", email))

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.userStorage.SaveUser(ctx, email, passHash)
	if errors.Is(err, models.ErrUserExists) {
		return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	if err != nil {
		log.Error("failed to save user", sl.Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

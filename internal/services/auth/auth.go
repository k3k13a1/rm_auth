package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/rizzmatch/rm_auth/internal/core/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
)

type Auth struct {
	log      *slog.Logger
	usrSaver UserSaver
	// usrProvider UserProvider
	tokenTTL time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid uuid.UUID, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, uid uuid.UUID) (bool, error)
}

func New(
	log *slog.Logger,
	userSaver UserSaver,
	// userProvider UserProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:      log,
		usrSaver: userSaver,
		// usrProvider: userProvider,
		tokenTTL: tokenTTL,
	}
}

// func (a *Auth) Login(ctx context.Context, email, password string) (string, error) {
// 	const op = "auth.Login"

// 	log := a.log.With(
// 		slog.String("op", op),
// 		slog.String("username", email),
// 	)

// 	log.Info("attempting to login user")

// 	user, err := a.usrProvider.User(ctx, email)
// 	if err != nil {
// 		if errors.Is(err, storage.ErrUserNotFound) {
// 			a.log.Warn("user not found", slog.String("error", err.Error()))

// 			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
// 		}

// 		a.log.Error("failed to user", slog.String("error", err.Error()))

// 		return "", fmt.Errorf("%s: %w", op, err)
// 	}

// 	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
// 		a.log.Info("invalid credentials", slog.String("error", err.Error()))

// 		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
// 	}

// 	log.Info("successfully logged in user")

// 	token, err := jwt.NewToken(user, a.tokenTTL)
// 	if err != nil {
// 		a.log.Error("failed to generate token", slog.String("error", err.Error()))
// 	}

// 	return token, nil
// }

func (a *Auth) Register(ctx context.Context, email string, pass string) (uuid.UUID, error) {
	const op = "auth.Register"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", email),
	)

	log.Info("attempting to register user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash password", slog.String("error", err.Error()))

		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.usrSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		log.Error("failed to save user", slog.String("error", err.Error()))

		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
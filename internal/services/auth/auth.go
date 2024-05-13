package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/rizzmatch/rm_auth/internal/core/jwt"
	"github.com/rizzmatch/rm_auth/internal/core/models"
	"github.com/rizzmatch/rm_auth/internal/storage"
	"github.com/rizzmatch/rm_auth/internal/storage/postgres"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
)

type Auth struct {
	usrSaver      UserSaver
	usrProvider   UserProvider
	phoneProvider PhoneProvider
	emailProvider EmailProvider
	tokenTTL      time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
	) (uid int, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, email string) (bool, error)
}

type PhoneProvider interface {
	AddPhone(ctx context.Context, id int64, phone string) (err error)
	EditPhone(ctx context.Context, id int64, phone string) (err error)
}

type EmailProvider interface {
	AddEmail(ctx context.Context, id int64, email string) (err error)
	EditEmail(ctx context.Context, id int64, email string) (err error)
}

func New(
	storage *postgres.Storage,
	// userSaver UserSaver,
	// userProvider UserProvider,
	// phoneProvider PhoneProvider,
	// emailProvider EmailProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		usrSaver:      storage,
		usrProvider:   storage,
		phoneProvider: storage,
		emailProvider: storage,
		tokenTTL:      tokenTTL,
	}
}

func (a *Auth) Login(ctx context.Context, email, password string) (string, error) {
	const op = "auth.Login"

	log := slog.With(
		slog.String("op", op),
		slog.String("username", email),
	)

	log.Info("attempting to login user")

	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", slog.String("error", err.Error()))

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("failed to user", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Info("invalid credentials", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("successfully logged in user")

	token, err := jwt.NewToken(user, a.tokenTTL)
	if err != nil {
		log.Error("failed to generate token", slog.String("error", err.Error()))
	}

	return token, nil
}

func (a *Auth) Register(ctx context.Context, email string, pass string) (int, error) {
	const op = "auth.Register"

	log := slog.With(
		slog.String("op", op),
		slog.String("username", email),
	)

	log.Info("attempting to register user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash password", slog.String("error", err.Error()))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.usrSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		log.Error("failed to save user", slog.String("error", err.Error()))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (a *Auth) IsAdmin(ctx context.Context, email string) (bool, error) {
	return a.usrProvider.IsAdmin(ctx, email)
}

func (a *Auth) AddPhone(ctx context.Context, id int64, phone string) error {
	return a.phoneProvider.AddPhone(ctx, id, phone)
}

func (a *Auth) EditPhone(ctx context.Context, id int64, phone string) error {
	return a.phoneProvider.EditPhone(ctx, id, phone)
}

func (a *Auth) AddEmail(ctx context.Context, id int64, email string) error {
	return a.emailProvider.AddEmail(ctx, id, email)
}

func (a *Auth) EditEmail(ctx context.Context, id int64, email string) error {
	return a.emailProvider.EditEmail(ctx, id, email)
}

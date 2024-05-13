package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

const (
	envLocal = "local"
	envProd  = "production"
)

type HandlerMiddleware struct {
	next slog.Handler
}

func NewHandlerMiddleware(next slog.Handler) *HandlerMiddleware {
	return &HandlerMiddleware{next: next}
}

func (h *HandlerMiddleware) Enabled(ctx context.Context, rec slog.Level) bool {
	return h.next.Enabled(ctx, rec)
}

func (h *HandlerMiddleware) Handle(ctx context.Context, rec slog.Record) error {
	if c, ok := ctx.Value(key).(logCtx); ok {
		if c.UserID != 0 {
			rec.Add("userID", c.UserID)
		}
		if c.Email != "" {
			rec.Add("email", c.Email)
		}
		if c.Phone != "" {
			rec.Add("phone", c.Phone)
		}
	}
	return h.next.Handle(ctx, rec)
}

func (h *HandlerMiddleware) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &HandlerMiddleware{next: h.next.WithAttrs(attrs)}
}

func (h *HandlerMiddleware) WithGroup(name string) slog.Handler {
	return &HandlerMiddleware{next: h.next.WithGroup(name)}
}

type logCtx struct {
	UserID int
	Email  string
	Phone  string
}

type keyType int

const key = keyType(0)

func WithLogUserID(ctx context.Context, userID int) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.UserID = userID
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{UserID: userID})
}

func WithLogEmail(ctx context.Context, email string) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.Email = email
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{Email: email})
}

func WithLogPhone(ctx context.Context, phone string) context.Context {
	if len(phone) > 5 {
		phone = string(phone[0]) + strings.Repeat("*", len(phone)-4) + phone[len(phone)-4:]
	}
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.Phone = phone
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{Phone: phone})
}

type errorWithLogCtx struct {
	next error
	ctx  logCtx
}

func (e *errorWithLogCtx) Error() string {
	return e.next.Error()
}

func WrapError(ctx context.Context, err error) error {
	c := logCtx{}
	if x, ok := ctx.Value(key).(logCtx); ok {
		c = x
	}
	return &errorWithLogCtx{
		next: err,
		ctx:  c,
	}
}

func ErrorCtx(ctx context.Context, err error) context.Context {
	if e, ok := err.(*errorWithLogCtx); ok { // в реальной жизни используйте error.As
		return context.WithValue(ctx, key, e.ctx)
	}
	return ctx
}

func SetupLogger(env string) {
	var handler slog.Handler

	switch env {
	case envLocal:
		handler = slog.Handler(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}),
		)
	case envProd:
		handler = slog.Handler(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			}),
		)
	default:
		panic("unknown environment")
	}

	handler = NewHandlerMiddleware(handler)

	slog.SetDefault(slog.New(handler))
}

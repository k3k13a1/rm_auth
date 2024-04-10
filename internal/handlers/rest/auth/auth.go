package authrest

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Auth interface {
	Register(
		ctx context.Context,
		email string,
		pass string,
	) (uuid.UUID, error)
}

// type serverAPI struct {
// 	auth Auth
// }

func Register(c echo.Context, a Auth) error {

	email := c.FormValue("email")
	pass := c.FormValue("password")

	uid, err := a.Register(context.TODO(), email, pass)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, uid)
}

package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	authorizationHeaderKey = "authorization"
	authorizationBearer    = "bearer"
	authorizationPayload   = "authorization_payload"
)

var ErrUnauthorizedToken = errors.New("invalid token")

var CustomLoggerConfig = middleware.LoggerConfig{
	Format: `{"severity":"${level}", "time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}",` +
		`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
		`"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
		`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
}

func (server *Server) authMiddleware(e echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		authorizationHeader := ctx.Request().Header.Get(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			return ctx.JSON(http.StatusUnauthorized, ErrUnauthorizedToken)
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			return ctx.JSON(http.StatusUnauthorized, ErrUnauthorizedToken)
		}

		if strings.ToLower(fields[0]) != authorizationBearer {
			return ctx.JSON(http.StatusUnauthorized, ErrUnauthorizedToken)
		}

		accessToken := fields[1]
		payload, err := server.token.VerifyToken(accessToken)
		if err != nil {
			return ctx.JSON(http.StatusUnauthorized, ErrUnauthorizedToken)
		}

		ctx.Set(authorizationPayload, payload)

		return nil
	}
}

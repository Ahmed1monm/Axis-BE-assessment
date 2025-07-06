package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

func RequestLogger(log zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			req := c.Request()
			res := c.Response()

			log.Info().
				Str("method", req.Method).
				Str("uri", req.RequestURI).
				Int("status", res.Status).
				Dur("duration", time.Since(start)).
				Msg("Request handled")

			return err
		}
	}
}

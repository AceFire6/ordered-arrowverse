package healthcheck

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/AceFire6/ordered-arrowverse/internal/build"
	"github.com/AceFire6/ordered-arrowverse/internal/customctx"
)

type Service interface {
	CheckCount() int
	RunHealthChecks(ctx context.Context) map[string]error
}

type Handler struct {
	healthCheckService Service
}

func NewHealthCheckHandler(healthCheckService Service) *Handler {
	return &Handler{
		healthCheckService: healthCheckService,
	}
}

func (hcHandler *Handler) HealthCheck(ctx echo.Context) error {
	cc := ctx.(*customctx.Context)

	checkCount := hcHandler.healthCheckService.CheckCount()
	// we give each check an average of 500ms to finish
	timeout := time.Duration(checkCount) * 500 * time.Millisecond

	timeoutCtx, cancel := context.WithTimeout(ctx.Request().Context(), timeout)
	defer cancel()

	// Change this to 500 if any of them failed
	statusCode := http.StatusOK
	healthCheckResponse := map[string]string{
		"api": "ok",
	}

	healthCheckResults := hcHandler.healthCheckService.RunHealthChecks(timeoutCtx)
	for name, checkErr := range healthCheckResults {
		if _, found := healthCheckResponse[name]; found {
			cc.Log.Warn().Str("check_name", name).Msg("Duplicate health check - skipping")

			continue
		}

		healthCheckResponse[name] = "ok"
		if checkErr != nil {
			// Only show the actual error if we're running in development
			if cc.Environment == build.Development {
				healthCheckResponse[name] = checkErr.Error()
			}

			healthCheckResponse[name] = "error"
		}
	}

	return ctx.JSON(statusCode, healthCheckResponse)
}

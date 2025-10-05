package healthcheck

import "context"

type healthChecker interface {
	RunCheck(ctx context.Context) error
	Name() string
}

type healthCheckService struct {
	healthChecks []healthChecker
}

func NewService(checkers ...healthChecker) *healthCheckService {
	return &healthCheckService{
		healthChecks: checkers,
	}
}

func (svc *healthCheckService) CheckCount() int {
	return len(svc.healthChecks)
}

func (svc *healthCheckService) RunHealthChecks(ctx context.Context) map[string]error {
	healthCheckResults := make(map[string]error)

	for _, check := range svc.healthChecks {
		healthCheckResults[check.Name()] = check.RunCheck(ctx)
	}

	return healthCheckResults
}

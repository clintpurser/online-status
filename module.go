package onlinestatus

import (
	"context"
	"fmt"
	"net/http"
	"time"

	sensor "go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

var (
	OnlineStatus = resource.NewModel("cdp", "online-status", "online-status")
)

func init() {
	resource.RegisterComponent(sensor.API, OnlineStatus,
		resource.Registration[sensor.Sensor, resource.NoNativeConfig]{
			Constructor: newOnlineStatusOnlineStatus,
		},
	)
}

type onlineStatusOnlineStatus struct {
	resource.AlwaysRebuild

	name resource.Name

	logger logging.Logger

	cancelCtx  context.Context
	cancelFunc func()
}

func newOnlineStatusOnlineStatus(ctx context.Context, deps resource.Dependencies, rawConf resource.Config, logger logging.Logger) (sensor.Sensor, error) {

	return NewOnlineStatus(ctx, deps, rawConf.ResourceName(), resource.NoNativeConfig{}, logger)

}

func NewOnlineStatus(ctx context.Context, deps resource.Dependencies, name resource.Name, _ resource.NoNativeConfig, logger logging.Logger) (sensor.Sensor, error) {

	cancelCtx, cancelFunc := context.WithCancel(context.Background())

	s := &onlineStatusOnlineStatus{
		name:       name,
		logger:     logger,
		cancelCtx:  cancelCtx,
		cancelFunc: cancelFunc,
	}
	return s, nil
}

func (s *onlineStatusOnlineStatus) Name() resource.Name {
	return s.name
}

func (s *onlineStatusOnlineStatus) Readings(ctx context.Context, extra map[string]interface{}) (map[string]interface{}, error) {
	// Create a custom HTTP client with a timeout, similar to the Python example.
	client := &http.Client{
		Timeout: 2 * time.Second, // 2-second timeout
	}

	req, err := http.NewRequestWithContext(ctx, "GET", "https://app.viam.com", nil)
	if err != nil {
		s.logger.Error("Failed to create HTTP request", "error", err)
		return map[string]interface{}{"online": 0}, nil // Treat request creation error as offline
	}

	resp, err := client.Do(req)
	if err != nil {
		s.logger.Debug("Failed to reach app.viam.com, treating as offline", "error", err)
		return map[string]interface{}{"online": 0}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return map[string]interface{}{"online": 1}, nil
	} else {
		s.logger.Debug("Received non-200 response code from app.viam.com", "statusCode", resp.StatusCode)
		return map[string]interface{}{"online": 0}, nil
	}
}

func (s *onlineStatusOnlineStatus) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *onlineStatusOnlineStatus) Close(context.Context) error {
	// Put close code here
	s.cancelFunc()
	return nil
}

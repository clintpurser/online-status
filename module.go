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

type Config struct {
	/*
		Put config attributes here. There should be public/exported fields
		with a `json` parameter at the end of each attribute.

		Example config struct:
			type Config struct {
				Pin   string `json:"pin"`
				Board string `json:"board"`
				MinDeg *float64 `json:"min_angle_deg,omitempty"`
			}

		If your model does not need a config, replace *Config in the init
		function with resource.NoNativeConfig
	*/
}

// Validate ensures all parts of the config are valid and important fields exist.
// Returns implicit required (first return) and optional (second return) dependencies based on the config.
// The path is the JSON path in your robot's config (not the `Config` struct) to the
// resource being validated; e.g. "components.0".
func (cfg *Config) Validate(path string) ([]string, []string, error) {
	// Add config validation code here
	return nil, nil, nil
}

type onlineStatusOnlineStatus struct {
	resource.AlwaysRebuild

	name resource.Name

	logger logging.Logger
	cfg    *Config

	cancelCtx  context.Context
	cancelFunc func()
}

func newOnlineStatusOnlineStatus(ctx context.Context, deps resource.Dependencies, rawConf resource.Config, logger logging.Logger) (sensor.Sensor, error) {
	conf, err := resource.NativeConfig[*Config](rawConf)
	if err != nil {
		return nil, err
	}

	return NewOnlineStatus(ctx, deps, rawConf.ResourceName(), conf, logger)

}

func NewOnlineStatus(ctx context.Context, deps resource.Dependencies, name resource.Name, conf *Config, logger logging.Logger) (sensor.Sensor, error) {

	cancelCtx, cancelFunc := context.WithCancel(context.Background())

	s := &onlineStatusOnlineStatus{
		name:       name,
		logger:     logger,
		cfg:        conf,
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

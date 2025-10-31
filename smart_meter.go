package ecoflow

import (
	"context"
)

// SmartMeter represents a smart meter device
// Ecoflow documentation: https://developer-eu.ecoflow.com/us/document/shp

type SmartMeter struct {
	c  *Client
	sn string
}

func (s *SmartMeter) GetSn() string {
	return s.sn
}

// GetAllParameters Get all parameters from the SmartMeter device
func (s *SmartMeter) GetAllParameters(ctx context.Context) (map[string]any, error) {
	return s.c.GetDeviceAllParameters(ctx, s.sn)
}

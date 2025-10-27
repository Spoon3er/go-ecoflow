package ecoflow

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Stream represents an Ecoflow Stream device
// API Documentation: https://developer-eu.ecoflow.com/us/document/bkw
type Stream struct {
	c  *Client
	sn string
}

func (s *Stream) GetSn() string {
	return s.sn
}

// Constants for Stream device commands
const (
	cmdId   = 17
	cmdFunc = 254
	dirDest = 1
	dirSrc  = 1
	dest    = 2
	needAck = true
)

type SwitchDevice string

const (
	AC1 SwitchDevice = "cfgRelay2Onoff"
	AC2 SwitchDevice = "cfgRelay3Onoff"
)

type operatingMode string

const (
	SelfPoweredMode operatingMode = "operateSelfPoweredOpen"
	AIMode          operatingMode = "operateIntelligentScheduleModeOpen"
)

// SetPowerSocket Devicename (AC1/AC2) on/off
// {"sn": "BKW2000000000001", "CmdId": 17, "CmdFunc": 254, "DirDest": 1, "DirSrc": 1, "Dest": 2, "NeedAck": true, "params": {"cfgRelay2Onoff": true}}
func (s *Stream) SetPowerSocket(ctx context.Context, switchDevice SwitchDevice, action bool) (*CmdSetResponse, error) {
	params := make(map[string]any)
	params[string(switchDevice)] = action
	return s.setParameter(ctx, params)
}

// SetBackupReserveLevel Set backup reserve level (3-95)
// {"sn": "BKW200000000000001", "CmdId": 17, "CmdFunc": 254, "DirDest": 1, "DirSrc": 1, "Dest": 2, "NeedAck": true, "params": {"backupReverseSoc": 20}}
func (s *Stream) SetBackupReserveLevel(ctx context.Context, lvl int) (*CmdSetResponse, error) {
	if lvl < 3 || lvl > 95 {
		return nil, fmt.Errorf("lvl must be between 3 and 95, got %d", lvl)
	}
	params := make(map[string]any)
	params["backupReverseSoc"] = lvl
	return s.setParameter(ctx, params)
}

// SetChargeLimit Set charge limit (50-95)
// {"sn": "BKW2000000000000001", "CmdId": 17, "CmdFunc": 254, "DirDest": 1, "DirSrc": 1, "Dest": 2, "NeedAck": true, "params": {"cmsMaxChgSoc": 80}}
func (s *Stream) SetChargeLimit(ctx context.Context, limit int) (*CmdSetResponse, error) {
	if limit < 50 || limit > 95 {
		return nil, fmt.Errorf("limit must be between 50 and 95, got %d", limit)
	}
	params := make(map[string]any)
	params["cmsMaxChgSoc"] = limit
	return s.setParameter(ctx, params)
}

// SetDisChargeLimit Set charge limit (3-50)
// {"sn": "BKW2000000000000001", "CmdId": 17, "CmdFunc": 254, "DirDest": 1, "DirSrc": 1, "Dest": 2, "NeedAck": true, "params": {"cmsMaxDisChgSoc": 80}}
func (s *Stream) SetDisChargeLimit(ctx context.Context, limit int) (*CmdSetResponse, error) {
	if limit < 3 || limit > 50 {
		return nil, fmt.Errorf("limit must be between 3 and 50, got %d", limit)
	}
	params := make(map[string]any)
	params["cmsMinDsgSoc"] = limit
	return s.setParameter(ctx, params)
}

// SetOperatingMode Set operating mode (SelfPoweredMode/AIMode)
// {"sn": "BKW2000000000001", "CmdId": 17, "CmdFunc": 254, "DirDest": 1, "DirSrc": 1, "Dest": 2, "NeedAck": true, "params": {"energyStrategyOperateMode": {"SelfPoweredMode": true, "AIMode": false}}}
func (s *Stream) SetOperatingMode(ctx context.Context, mode operatingMode, value bool) (*CmdSetResponse, error) {
	var selfPowered, aiMode bool

	if mode == SelfPoweredMode {
		selfPowered = value
		aiMode = !value
	} else {
		aiMode = value
		selfPowered = !value
	}

	params := map[string]any{
		"cfgEnergyStrategyOperateMode": map[string]bool{
			string(SelfPoweredMode): selfPowered,
			string(AIMode):          aiMode,
		},
	}

	return s.setParameter(ctx, params)
}

// SetFeedInControl Set feed-in control (on/off)
// {"sn": "BKW2000000000001", "CmdId": 17, "CmdFunc": 254, "DirDest": 1, "DirSrc": 1, "Dest": 2, "NeedAck": true, "params": {"cfgFeedGridMode": 2}}
func (s *Stream) SetFeedInControl(ctx context.Context, value string) (*CmdSetResponse, error) {
	params := make(map[string]any)

	switch value {
	case "off":
		params["cfgFeedGridMode"] = 1
	case "on":
		params["cfgFeedGridMode"] = 2
	default:
		return nil, fmt.Errorf("value must be 'on' or 'off', got %s", value)
	}

	return s.setParameter(ctx, params)
}

// GetEnergyIndependence
// {"sn":"BK2000000000001","params":{"beginTime":"2023-10-01 00:00:00","endTime":"2023-10-10 23:59:59","code":"BK621-App-HOME-INDEPENDENCE-PERCENT-FLOW-indep-progress_bar-NOTDISTINGUISH-MASTER_DATA"}}
func (s *Stream) GetEnergyIndependence(ctx context.Context, beginTime, endTime time.Time) (*GetHistoricalDataResponse, error) {
	code := "BK621-App-HOME-INDEPENDENCE-PERCENT-FLOW-indep-progress_bar-NOTDISTINGUISH-MASTER_DATA"
	return s.getHistoricalData(ctx, beginTime, endTime, code)
}

// GetEnvironmentalImpact
// {"sn":"BK2000000000001","params":{"beginTime":"2023-10-01 00:00:00","endTime":"2023-10-10 23:59:59","code":"BK621-App-HOME-SAVING-CURRENCY-FLOW-earnings-progress_arc-NOTDISTINGUISH-MASTER_DATA"}}
func (s *Stream) GetEnvironmentalImpact(ctx context.Context, beginTime, endTime time.Time) (*GetHistoricalDataResponse, error) {
	code := "BK621-App-HOME-CO2-WEIGHT-FLOW-impact-progress_arc-NOTDISTINGUISH-MASTER_DATA"
	return s.getHistoricalData(ctx, beginTime, endTime, code)
}

// GetTotalEnergySavings
// {"sn":"BK2000000000001","params":{"beginTime":"2023-10-01 00:00:00","endTime":"2023-10-10 23:59:59","code":"BK621-App-HOME-SOLAR-ENERGY-FLOW-solor-line-NOTDISTINGUISH-MASTER_DATA"}}
func (s *Stream) GetTotalEnergySavings(ctx context.Context, beginTime, endTime time.Time) (*GetHistoricalDataResponse, error) {
	code := "BK621-App-HOME-SOLAR-ENERGY-FLOW-solor-line-NOTDISTINGUISH-MASTER_DATA"
	return s.getHistoricalData(ctx, beginTime, endTime, code)
}

// GetElectricityConsumption
// {"sn":"BK2000000000001","params":{"beginTime":"2023-10-01 00:00:00","endTime":"2023-10-10 23:59:59","code":"BK621-App-HOME-LOAD-ENERGY-FLOW-consumption-prop_arc-NOTDISTINGUISH-MASTER_DATA"}}
func (s *Stream) GetElectricityConsumption(ctx context.Context, beginTime, endTime time.Time) (*GetHistoricalDataResponse, error) {
	code := "BK621-App-HOME-LOAD-ENERGY-FLOW-consumption-prop_arc-NOTDISTINGUISH-MASTER_DATA"
	return s.getHistoricalData(ctx, beginTime, endTime, code)
}

// GetGrid
// {"sn":"BK2000000000001","params":{"beginTime":"2023-10-01 00:00:00","endTime":"2023-10-10 23:59:59","code":"BK621-App-HOME-GRID-ENERGY-FLOW-grid_prop_bar-NOTDISTINGUISH-MASTER_DATA"}}
func (s *Stream) GetGrid(ctx context.Context, beginTime, endTime time.Time) (*GetHistoricalDataResponse, error) {
	code := "BK621-App-HOME-GRID-ENERGY-FLOW-grid_prop_bar-NOTDISTINGUISH-MASTER_DATA"
	return s.getHistoricalData(ctx, beginTime, endTime, code)
}

// GetBatteryCharging/DischargingPower
// {"sn":"BK2000000000001","params":{"beginTime":"2023-10-01 00:00:00","endTime":"2023-10-10 23:59:59","code":"BK621-App-HOME-SOC-ENERGY-FLOW-battery-prop_bar-NOTDISTINGUISH-MASTER_DATA"}}
func (s *Stream) GetBatteryChargingDischargingPower(ctx context.Context, beginTime, endTime time.Time) (*GetHistoricalDataResponse, error) {
	code := "BK621-App-HOME-SOC-ENERGY-FLOW-battery-prop_bar-NOTDISTINGUISH-MASTER_DATA"
	return s.getHistoricalData(ctx, beginTime, endTime, code)
}

// getHistoricalData internal helper function to get historical data from the Stream device
func (s *Stream) getHistoricalData(ctx context.Context, beginTime, endTime time.Time, code string) (*GetHistoricalDataResponse, error) {
	params := GetHistoricalDataParams{}
	params.BeginTime = beginTime.Format("2006-01-02 15:04:05")
	params.EndTime = endTime.Format("2006-01-02 15:04:05")
	params.Code = code

	return s.c.GetDeviceHistoricalData(ctx, s.sn, &params)
}

// GetParameter Get specific parameters from the Stream device
func (s *Stream) GetParameter(ctx context.Context, params []string) (*GetCmdResponse, error) {
	return s.c.GetDeviceParameters(ctx, s.sn, params)
}

// GetAllParameters Get all parameters from the Stream device
func (s *Stream) GetAllParameters(ctx context.Context) (map[string]any, error) {
	return s.c.GetDeviceAllParameters(ctx, s.sn)
}

// setParameter internal helper function to set device parameters
func (s *Stream) setParameter(ctx context.Context, params map[string]any) (*CmdSetResponse, error) {
	cmdReq := CmdUltraSetRequest{
		Id:      fmt.Sprint(time.Now().UnixMilli()),
		CmdId:   cmdId,
		CmdFunc: cmdFunc,
		DirDest: dirDest,
		DirSrc:  dirSrc,
		Dest:    dest,
		NeedAck: needAck,
		Sn:      s.sn,
		Params:  params,
	}

	jsonData, err := json.Marshal(cmdReq)
	if err != nil {
		return nil, err
	}

	var req map[string]any
	err = json.Unmarshal(jsonData, &req)
	if err != nil {
		return nil, err
	}

	return s.c.SetDeviceParameter(ctx, req)
}

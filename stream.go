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

	request, err := BuildRequest(s.sn, params)
	if err != nil {
		return nil, err
	}
	return s.c.SetDeviceParameter(ctx, request)
}

// SetBackupReserveLevel Set backup reserve level (3-95)
// {"sn": "BKW200000000000001", "CmdId": 17, "CmdFunc": 254, "DirDest": 1, "DirSrc": 1, "Dest": 2, "NeedAck": true, "params": {"backupReverseSoc": 20}}
func (s *Stream) SetBackupReserveLevel(ctx context.Context, lvl int) (*CmdSetResponse, error) {
	if lvl < 3 || lvl > 95 {
		return nil, fmt.Errorf("lvl must be between 3 and 95, got %d", lvl)
	}
	params := make(map[string]any)
	params["cfgBackupReverseSoc"] = lvl

	request, err := BuildRequest(s.sn, params)
	if err != nil {
		return nil, err
	}
	return s.c.SetDeviceParameter(ctx, request)
}

// SetChargeLimit Set charge limit (50-95)
// {"sn": "BKW2000000000000001", "CmdId": 17, "CmdFunc": 254, "DirDest": 1, "DirSrc": 1, "Dest": 2, "NeedAck": true, "params": {"cmsMaxChgSoc": 80}}
func (s *Stream) SetChargeLimit(ctx context.Context, limit int) (*CmdSetResponse, error) {
	if limit < 50 || limit > 95 {
		return nil, fmt.Errorf("limit must be between 50 and 95, got %d", limit)
	}
	params := make(map[string]any)
	params["cmsMaxChgSoc"] = limit

	request, err := BuildRequest(s.sn, params)
	if err != nil {
		return nil, err
	}
	return s.c.SetDeviceParameter(ctx, request)
}

// SetDisChargeLimit Set charge limit (3-50)
// {"sn": "BKW2000000000000001", "CmdId": 17, "CmdFunc": 254, "DirDest": 1, "DirSrc": 1, "Dest": 2, "NeedAck": true, "params": {"cmsMaxDisChgSoc": 80}}
func (s *Stream) SetDisChargeLimit(ctx context.Context, limit int) (*CmdSetResponse, error) {
	if limit < 3 || limit > 50 {
		return nil, fmt.Errorf("limit must be between 3 and 50, got %d", limit)
	}
	params := make(map[string]any)
	params["cmsMinDsgSoc"] = limit

	request, err := BuildRequest(s.sn, params)
	if err != nil {
		return nil, err
	}
	return s.c.SetDeviceParameter(ctx, request)
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

	request, err := BuildRequest(s.sn, params)
	if err != nil {
		return nil, err
	}
	return s.c.SetDeviceParameter(ctx, request)
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

	request, err := BuildRequest(s.sn, params)
	if err != nil {
		return nil, err
	}
	return s.c.SetDeviceParameter(ctx, request)
}

// GetEnergyIndependence retrieves energy independence data for the specified time range
// {"sn":"BK2000000000001","params":{"beginTime":"2023-10-01 00:00:00","endTime":"2023-10-10 23:59:59","code":"BK621-App-HOME-INDEPENDENCE-PERCENT-FLOW-indep-progress_bar-NOTDISTINGUISH-MASTER_DATA"}}
func (s *Stream) GetEnergyIndependence(ctx context.Context, beginTime, endTime time.Time) (*HistoricalDataResponse, error) {
	code := "BK621-App-HOME-INDEPENDENCE-PERCENT-FLOW-indep-progress_bar-NOTDISTINGUISH-MASTER_DATA"
	return s.getHistoricalData(ctx, beginTime, endTime, code)
}

// GetEnvironmentalImpact retrieves environmental impact data for the specified time range
// {"sn":"BK2000000000001","params":{"beginTime":"2023-10-01 00:00:00","endTime":"2023-10-10 23:59:59","code":"BK621-App-HOME-SAVING-CURRENCY-FLOW-earnings-progress_arc-NOTDISTINGUISH-MASTER_DATA"}}
func (s *Stream) GetEnvironmentalImpact(ctx context.Context, beginTime, endTime time.Time) (*HistoricalDataResponse, error) {
	code := "BK621-App-HOME-CO2-WEIGHT-FLOW-impact-progress_arc-NOTDISTINGUISH-MASTER_DATA"
	return s.getHistoricalData(ctx, beginTime, endTime, code)
}

// GetTotalEnergySavings retrieves total energy savings data for the specified time range
// {"sn":"BK2000000000001","params":{"beginTime":"2023-10-01 00:00:00","endTime":"2023-10-10 23:59:59","code":"BK621-App-HOME-SOLAR-ENERGY-FLOW-solor-line-NOTDISTINGUISH-MASTER_DATA"}}
func (s *Stream) GetTotalEnergySavings(ctx context.Context, beginTime, endTime time.Time) (*HistoricalDataResponse, error) {
	code := "BK621-App-HOME-SOLAR-ENERGY-FLOW-solor-line-NOTDISTINGUISH-MASTER_DATA"
	return s.getHistoricalData(ctx, beginTime, endTime, code)
}

// GetElectricityConsumption retrieves electricity consumption data for the specified time range
// {"sn":"BK2000000000001","params":{"beginTime":"2023-10-01 00:00:00","endTime":"2023-10-10 23:59:59","code":"BK621-App-HOME-LOAD-ENERGY-FLOW-consumption-prop_arc-NOTDISTINGUISH-MASTER_DATA"}}
func (s *Stream) GetElectricityConsumption(ctx context.Context, beginTime, endTime time.Time) (*HistoricalDataResponse, error) {
	code := "BK621-App-HOME-LOAD-ENERGY-FLOW-consumption-prop_arc-NOTDISTINGUISH-MASTER_DATA"
	return s.getHistoricalData(ctx, beginTime, endTime, code)
}

// GetGrid retrieves grid energy data for the specified time range
// {"sn":"BK2000000000001","params":{"beginTime":"2023-10-01 00:00:00","endTime":"2023-10-10 23:59:59","code":"BK621-App-HOME-GRID-ENERGY-FLOW-grid_prop_bar-NOTDISTINGUISH-MASTER_DATA"}}
func (s *Stream) GetGrid(ctx context.Context, beginTime, endTime time.Time) (*HistoricalDataResponse, error) {
	code := "BK621-App-HOME-GRID-ENERGY-FLOW-grid_prop_bar-NOTDISTINGUISH-MASTER_DATA"
	return s.getHistoricalData(ctx, beginTime, endTime, code)
}

// GetBatteryChargingDischargingPower retrieves battery charging/discharging power data for the specified time range
// {"sn":"BK2000000000001","params":{"beginTime":"2023-10-01 00:00:00","endTime":"2023-10-10 23:59:59","code":"BK621-App-HOME-SOC-ENERGY-FLOW-battery-prop_bar-NOTDISTINGUISH-MASTER_DATA"}}
func (s *Stream) GetBatteryChargingDischargingPower(ctx context.Context, beginTime, endTime time.Time) (*HistoricalDataResponse, error) {
	code := "BK621-App-HOME-SOC-ENERGY-FLOW-battery-prop_bar-NOTDISTINGUISH-MASTER_DATA"
	return s.getHistoricalData(ctx, beginTime, endTime, code)
}

// getHistoricalData internal helper function to get historical data from the Stream device
func (s *Stream) getHistoricalData(ctx context.Context, beginTime, endTime time.Time, code string) (*HistoricalDataResponse, error) {
	params := HistoricalDataParams{}
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

// BuildRequest internal helper function to set device parameters
func BuildRequest(sn string, params map[string]any) (map[string]any, error) {
	req := map[string]any{
		"id":      int(time.Now().UnixMilli()),
		"version": "1.0",
		"cmdId":   cmdId,
		"cmdFunc": cmdFunc,
		"dirDest": dirDest,
		"dirSrc":  dirSrc,
		"dest":    dest,
		"needAck": needAck,
		"sn":      sn,
		"params":  params,
	}

	return req, nil
}

// MQTT Quota Parameters for Stream devices
// These types represent the real-time data received from MQTT quota topic:
// /open/{certificateAccount}/{sn}/quota

// StreamPowerParams holds all power-related parameters for Stream devices
type StreamPowerParams struct {
	PowGetSysGrid             *float64 `json:"powGetSysGrid,omitempty"`
	PowGetSysLoad             *float64 `json:"powGetSysLoad,omitempty"`
	PowGetSysLoadFromGrid     *float64 `json:"powGetSysLoadFromGrid,omitempty"`
	PowGetSysLoadFromBp       *float64 `json:"powGetSysLoadFromBp,omitempty"`
	PowGetSysLoadFromPv       *float64 `json:"powGetSysLoadFromPv,omitempty"`
	PowGetBpCms               *float64 `json:"powGetBpCms,omitempty"`
	PowGetPv                  *float64 `json:"powGetPv,omitempty"`
	PowGetPv2                 *float64 `json:"powGetPv2,omitempty"`
	PowGetPv3                 *float64 `json:"powGetPv3,omitempty"`
	PowGetPv4                 *float64 `json:"powGetPv4,omitempty"`
	PowGetPvSum               *float64 `json:"powGetPvSum,omitempty"`
	PowGetSchuko1             *float64 `json:"powGetSchuko1,omitempty"`
	PowGetSchuko2             *float64 `json:"powGetSchuko2,omitempty"`
	PowSysAcInMax             *float64 `json:"powSysAcInMax,omitempty"`
	PowSysAcOutMax            *float64 `json:"powSysAcOutMax,omitempty"`
	SocketMeasurePower        *float64 `json:"socketMeasurePower,omitempty"`
	GridConnectionPower       *float64 `json:"gridConnectionPower,omitempty"`
	SysGridConnectionPower    *float64 `json:"sysGridConnectionPower,omitempty"`
	PowConsumptionMeasurement *float64 `json:"powConsumptionMeasurement,omitempty"`
}

// StreamBMSParams holds Battery Management System parameters for Stream devices
type StreamBMSParams struct {
	BmsBattSoc     *float64 `json:"bmsBattSoc,omitempty"`
	BmsBattSoh     *float64 `json:"bmsBattSoh,omitempty"`
	BmsDesignCap   *float64 `json:"bmsDesignCap,omitempty"`
	BmsBattHeating *bool    `json:"bmsBattHeating,omitempty"`
	BmsChgDsgState *float64 `json:"bmsChgDsgState,omitempty"`
	BmsChgRemTime  *float64 `json:"bmsChgRemTime,omitempty"`
	BmsDsgRemTime  *float64 `json:"bmsDsgRemTime,omitempty"`
	BmsMaxCellTemp *float64 `json:"bmsMaxCellTemp,omitempty"`
	BmsMinCellTemp *float64 `json:"bmsMinCellTemp,omitempty"`
	BmsMaxMosTemp  *float64 `json:"bmsMaxMosTemp,omitempty"`
	BmsMinMosTemp  *float64 `json:"bmsMinMosTemp,omitempty"`
}

// StreamCMSParams holds Central Management System parameters for Stream devices
type StreamCMSParams struct {
	CmsBattSoc        *float64 `json:"cmsBattSoc,omitempty"`
	CmsBattSoh        *float64 `json:"cmsBattSoh,omitempty"`
	CmsBattFullEnergy *float64 `json:"cmsBattFullEnergy,omitempty"`
	CmsBattPowInMax   *float64 `json:"cmsBattPowInMax,omitempty"`
	CmsBattPowOutMax  *float64 `json:"cmsBattPowOutMax,omitempty"`
	CmsBmsRunState    *float64 `json:"cmsBmsRunState,omitempty"`
	CmsChgDsgState    *float64 `json:"cmsChgDsgState,omitempty"`
	CmsChgRemTime     *float64 `json:"cmsChgRemTime,omitempty"`
	CmsDsgRemTime     *float64 `json:"cmsDsgRemTime,omitempty"`
	CmsMaxChgSoc      *float64 `json:"cmsMaxChgSoc,omitempty"`
	CmsMinDsgSoc      *float64 `json:"cmsMinDsgSoc,omitempty"`
}

// StreamPVParams holds Photovoltaic input parameters for Stream devices
type StreamPVParams struct {
	PlugInInfoPvVol   *float64 `json:"plugInInfoPvVol,omitempty"`
	PlugInInfoPvAmp   *float64 `json:"plugInInfoPvAmp,omitempty"`
	PlugInInfoPvFlag  *bool    `json:"plugInInfoPvFlag,omitempty"`
	PlugInInfoPv2Vol  *float64 `json:"plugInInfoPv2Vol,omitempty"`
	PlugInInfoPv2Amp  *float64 `json:"plugInInfoPv2Amp,omitempty"`
	PlugInInfoPv2Flag *bool    `json:"plugInInfoPv2Flag,omitempty"`
	PlugInInfoPv3Vol  *float64 `json:"plugInInfoPv3Vol,omitempty"`
	PlugInInfoPv3Amp  *float64 `json:"plugInInfoPv3Amp,omitempty"`
	PlugInInfoPv3Flag *bool    `json:"plugInInfoPv3Flag,omitempty"`
	PlugInInfoPv4Vol  *float64 `json:"plugInInfoPv4Vol,omitempty"`
	PlugInInfoPv4Amp  *float64 `json:"plugInInfoPv4Amp,omitempty"`
	PlugInInfoPv4Flag *bool    `json:"plugInInfoPv4Flag,omitempty"`
}

// StreamSystemParams holds system configuration parameters for Stream devices
type StreamSystemParams struct {
	ModuleWifiRssi *float64 `json:"moduleWifiRssi,omitempty"`
	Brightness     *float64 `json:"brightness,omitempty"`
	UtcTimezone    *float64 `json:"utcTimezone,omitempty"`
	UtcTimezoneId  *string  `json:"utcTimezoneId,omitempty"`
	DevCtrlStatus  *float64 `json:"devCtrlStatus,omitempty"`
}

// Helper functions to extract parameters from MQTT payloads
// Each function takes a JSON payload as input and returns the corresponding struct
func ExtractStreamPowerParams(payload []byte) (StreamPowerParams, error) {
	var params StreamPowerParams
	err := json.Unmarshal(payload, &params)
	if err != nil {
		return StreamPowerParams{}, err
	}
	return params, nil
}

func ExtractStreamBMSParams(payload []byte) (StreamBMSParams, error) {
	var params StreamBMSParams
	err := json.Unmarshal(payload, &params)
	if err != nil {
		return StreamBMSParams{}, err
	}
	return params, nil
}

func ExtractStreamCMSParams(payload []byte) (StreamCMSParams, error) {
	var params StreamCMSParams
	err := json.Unmarshal(payload, &params)
	if err != nil {
		return StreamCMSParams{}, err
	}
	return params, nil
}

func ExtractStreamPVParams(payload []byte) (StreamPVParams, error) {
	var params StreamPVParams
	err := json.Unmarshal(payload, &params)
	if err != nil {
		return StreamPVParams{}, err
	}
	return params, nil
}

func ExtractStreamSystemParams(payload []byte) (StreamSystemParams, error) {
	var params StreamSystemParams
	err := json.Unmarshal(payload, &params)
	if err != nil {
		return StreamSystemParams{}, err
	}
	return params, nil
}

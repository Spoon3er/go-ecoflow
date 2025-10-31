package ecoflow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

const (
	ecoflowApiUrl              = "https://api.ecoflow.com"
	deviceListUrl              = "/iot-open/sign/device/list"
	getAllQuoteUrl             = "/iot-open/sign/device/quota/all"
	setDeviceFunctionUrl       = "/iot-open/sign/device/quota"
	getDeviceFunctionUrl       = "/iot-open/sign/device/quota"
	getDeviceHistoricalDataUrl = "/iot-open/sign/device/quota/data"
)

type GridFrequency int

const (
	GridFrequency50Hz GridFrequency = 1
	GridFrequency60Hz GridFrequency = 2
)

type ModuleType int

const (
	ModuleTypePd       ModuleType = 1
	ModuleTypeBms      ModuleType = 2
	ModuleTypeInv      ModuleType = 3
	ModuleTypeBmsSlave ModuleType = 4
	ModuleTypeMppt     ModuleType = 5
)

type TemperatureUnit int

const (
	TemperatureUnitCelsius    TemperatureUnit = 0
	TemperatureUnitFahrenheit TemperatureUnit = 1
)

type Client struct {
	httpClient  *http.Client //can be customized if required
	accessToken string
	secretToken string
	baseUrl     string
}

// NewEcoflowClient with default http client
func NewEcoflowClient(accessToken, secretToken string, options ...func(*Client)) *Client {
	c := &Client{
		httpClient:  &http.Client{},
		accessToken: accessToken,
		secretToken: secretToken,
		baseUrl:     ecoflowApiUrl,
	}

	for _, o := range options {
		o(c)
	}
	return c
}

func WithBaseUrl(url string) func(client *Client) {
	return func(s *Client) {
		s.baseUrl = url
	}
}

func WithHttpClient(c *http.Client) func(client *Client) {
	return func(s *Client) {
		s.httpClient = c
	}
}

func (c *Client) GetPowerStation(sn string) *PowerStation {
	return &PowerStation{
		c:  c,
		sn: sn,
	}
}

func (c *Client) GetPowerStationPro(sn string) *PowerStationPro {
	return &PowerStationPro{
		c:  c,
		sn: sn,
	}
}

func (c *Client) GetPowerStreamMicroInverter(sn string) *PowerStreamMicroInverter {
	return &PowerStreamMicroInverter{
		c:  c,
		sn: sn,
	}
}

func (c *Client) GetSmartHomePanel(sn string) *SmartHomePanel {
	return &SmartHomePanel{
		c:  c,
		sn: sn,
	}
}

func (c *Client) GetSmartPlug(sn string) *SmartPlug {
	return &SmartPlug{
		c:  c,
		sn: sn,
	}
}

func (c *Client) GetWaveAirConditioner(sn string) *WaveAirConditioner {
	return &WaveAirConditioner{
		c:  c,
		sn: sn,
	}
}

func (c *Client) GetGlacier(sn string) *Glacier {
	return &Glacier{
		c:  c,
		sn: sn,
	}
}

func (c *Client) GetPowerKit(sn string, moduleSn string) *PowerKit {
	return &PowerKit{
		c:        c,
		sn:       sn,
		moduleSn: moduleSn,
	}
}

func (c *Client) GetStream(sn string) *Stream {
	return &Stream{
		c:  c,
		sn: sn,
	}
}

func (c *Client) GetSmartMeter(sn string) *SmartMeter {
	return &SmartMeter{
		c:  c,
		sn: sn,
	}
}

type SettingSwitcher int

const (
	SettingEnabled  SettingSwitcher = 1
	SettingDisabled SettingSwitcher = 0
)

type DeviceListResponse struct {
	Code            string       `json:"code"`
	Message         string       `json:"message"`
	Devices         []DeviceInfo `json:"data"`
	EagleEyeTraceID string       `json:"eagleEyeTraceId"`
	Tid             string       `json:"tid"`
}

type DeviceInfo struct {
	SN     string `json:"sn"`
	Online int    `json:"online"`
}

// GetDeviceList executes a request to get the list of devises linked to the user account. Shared devices are not included
// If the response parameter "code" is not 0, then there is an error. Error code and error message are returned
func (c *Client) GetDeviceList(ctx context.Context) (*DeviceListResponse, error) {
	request := NewHttpRequest(c.httpClient, "GET", c.baseUrl+deviceListUrl, nil, c.accessToken, c.secretToken)
	response, err := request.Execute(ctx)
	if err != nil {
		return nil, err
	}
	var deviceResponse DeviceListResponse

	err = json.Unmarshal(response, &deviceResponse)
	if err != nil {
		return nil, err
	}

	if deviceResponse.Code != "0" {
		return &deviceResponse, fmt.Errorf("can't get device list, error code: %s, error message: %s", deviceResponse.Code, deviceResponse.Message)
	}
	return &deviceResponse, nil
}

type CmdSetRequest struct {
	Id          string         `json:"id"`
	OperateType string         `json:"operateType,omitempty"`
	ModuleType  ModuleType     `json:"moduleType,omitempty"`
	CmdCode     string         `json:"cmdCode,omitempty"`
	Sn          string         `json:"sn"`
	Params      map[string]any `json:"params"`
}

type CmdUltraSetRequest struct {
	Id      string         `json:"id"`
	CmdId   int            `json:"cmdId"`
	CmdFunc int            `json:"cmdFunc"`
	DirDest int            `json:"dirDest"`
	DirSrc  int            `json:"dirSrc"`
	Dest    int            `json:"dest"`
	NeedAck bool           `json:"needAck"`
	Sn      string         `json:"sn"`
	Params  map[string]any `json:"params"`
}

type CmdSetResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func getParamsEnabled(enabled SettingSwitcher) map[string]any {
	params := make(map[string]any)
	params["enabled"] = enabled
	return params
}

// SetDeviceParameter exporter function to set device's settings.The request is a JSON map that will be sent to the server
// Each device has its own request structure so this function works for all types of devices.
// This function can be used even if your device type is not supported by this library
func (c *Client) SetDeviceParameter(ctx context.Context, request map[string]any) (*CmdSetResponse, error) {
	slog.Debug("SetDeviceParameter", "request", request)

	r := NewHttpRequest(c.httpClient, "PUT", c.baseUrl+setDeviceFunctionUrl, request, c.accessToken, c.secretToken)

	response, err := r.Execute(ctx)
	if err != nil {
		return nil, err
	}

	slog.Debug("SetDeviceParameter", "response", string(response))

	var cmdResponse *CmdSetResponse
	err = json.Unmarshal(response, &cmdResponse)
	if err != nil {
		return nil, err
	}

	return cmdResponse, nil
}

type GetCmdRequest struct {
	Sn     string        `json:"sn"`
	Params GetParamsList `json:"params"`
}

type GetParamsList struct {
	Quotas []string `json:"quotas"`
}

type GetCmdResponse struct {
	Code            string         `json:"code"`
	Message         string         `json:"message"`
	Data            map[string]any `json:"data"`
	EagleEyeTraceID string         `json:"eagleEyeTraceId"`
	Tid             string         `json:"tid"`
}

// GetDeviceParameters returns specified parameters for device
// This is a generic function that works for all types of devices
func (c *Client) GetDeviceParameters(ctx context.Context, deviceSN string, params []string) (*GetCmdResponse, error) {
	if len(params) == 0 {
		return nil, errors.New("parameters are mandatory")
	}
	if deviceSN == "" {
		return nil, errors.New("device SN is mandatory")
	}

	req := GetCmdRequest{
		Sn: deviceSN,
		Params: GetParamsList{
			Quotas: params,
		},
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	var reqParams map[string]any
	err = json.Unmarshal(jsonData, &reqParams)
	if err != nil {
		return nil, err
	}

	r := NewHttpRequest(c.httpClient, "POST", c.baseUrl+getDeviceFunctionUrl, reqParams, c.accessToken, c.secretToken)

	response, err := r.Execute(ctx)
	if err != nil {
		return nil, err
	}

	var getCmdResponse *GetCmdResponse

	err = json.Unmarshal(response, &getCmdResponse)
	if err != nil {
		return nil, err
	}

	if getCmdResponse.Code != "0" {
		return getCmdResponse, fmt.Errorf("can't get parameters, error code %s", getCmdResponse.Code)
	}

	return getCmdResponse, nil
}

// GetDeviceAllParameters executes a request to get the raw parameters ("as is") for a specific device.
// This function works for all types of devices.
// It returns a map[string]any containing the parameters and an error if any. The value type is mostly int, for some parameters it's float64 or []int
// If the response parameter "code" is not "0", then there is an error and the error message is returned.
// The parameters are taken from the Ecoflow response, "data" field
// If the response is not valid or cannot be processed, an error is returned.
func (c *Client) GetDeviceAllParameters(ctx context.Context, deviceSn string) (map[string]any, error) {
	requestParams := make(map[string]any)
	requestParams["sn"] = deviceSn

	request := NewHttpRequest(c.httpClient, "GET", c.baseUrl+getAllQuoteUrl, requestParams, c.accessToken, c.secretToken)
	response, err := request.Execute(ctx)

	if err != nil {
		return nil, err
	}

	var jsonData map[string]any
	err = json.Unmarshal(response, &jsonData)
	if err != nil {
		return nil, err
	}

	if code, ok := jsonData["code"].(string); !ok || code != "0" {
		return nil, fmt.Errorf("can't get parameters, error code %s", code)
	}

	dataMap, ok := jsonData["data"].(map[string]any)

	if !ok {
		return nil, errors.New("response is not valid, can't process it")
	}

	return dataMap, err
}

type GetHistoricalDataResponse struct {
	Code            string           `json:"code"`
	Message         string           `json:"message"`
	Data            []GetCmdResponse `json:"data"`
	EagleEyeTraceID string           `json:"eagleEyeTraceId"`
	Tid             string           `json:"tid"`
}

type GetHistoricalDataParams struct {
	BeginTime string `json:"beginTime"`
	EndTime   string `json:"endTime"`
	Code      string `json:"code"`
}

type GetHistoricalDataRequest struct {
	Sn     string                  `json:"sn"`
	Params GetHistoricalDataParams `json:"params"`
}

// GetDeviceHistoricalData retrieves historical data for a specific device based on provided parameters.
// It requires the device serial number and a list of parameters to fetch historical data for.
// If the response parameter "code" is not "0", then there is an error and the error message is returned.

func (c *Client) GetDeviceHistoricalData(ctx context.Context, deviceSn string, params *GetHistoricalDataParams) (*GetHistoricalDataResponse, error) {
	if params == nil {
		return nil, errors.New("parameters are mandatory")
	}
	if deviceSn == "" {
		return nil, errors.New("device SN is mandatory")
	}

	req := GetHistoricalDataRequest{
		Sn:     deviceSn,
		Params: *params,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	var reqParams map[string]any
	err = json.Unmarshal(jsonData, &reqParams)
	if err != nil {
		return nil, err
	}

	r := NewHttpRequest(c.httpClient, "POST", c.baseUrl+getDeviceHistoricalDataUrl, reqParams, c.accessToken, c.secretToken)

	response, err := r.Execute(ctx)
	if err != nil {
		return nil, err
	}

	var getHistoricalDataResponse *GetHistoricalDataResponse

	err = json.Unmarshal(response, &getHistoricalDataResponse)
	if err != nil {
		return nil, err
	}

	if getHistoricalDataResponse.Code != "0" {
		return getHistoricalDataResponse, fmt.Errorf("can't get historical data, error code %s", getHistoricalDataResponse.Code)
	}

	return getHistoricalDataResponse, nil

}

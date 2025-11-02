//go:build exclude

package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Spoon3er/go-ecoflow"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	accessKey := os.Getenv("ACCESS_KEY")
	secretKey := os.Getenv("SECRET_KEY")

	if accessKey == "" || secretKey == "" {
		slog.Error("AccessKey and SecretKey are mandatory")
		return
	}

	//create new client.
	client := ecoflow.NewEcoflowClient(accessKey, secretKey)

	// creating new client with options. Current supports two options:
	// 1. custom ecoflow base url (can be used with proxies, or if they change the url)
	// 2. custom http client

	// client = ecoflow.NewEcoflowClient(accessKey, secretKey,
	//	ecoflow.WithBaseUrl("https://ecoflow-api.example.com"),
	//	ecoflow.WithHttpClient(customHttpClient()),
	//)

	//get all linked ecoflow devices. Returns SN and online status
	client.GetDeviceList(context.Background())

	ctx := context.Background()

	// get set / get functions for power stations. PRO version is not currently implemented
	ps := client.GetStream(os.Getenv("DEVICE_SN"))

	// get functions
	params, err := ps.GetAllParameters(ctx)
	if err != nil {
		slog.Error("Failed to get parameters", "error", err)
	} else {
		slog.Info("Power Station All Parameters", "params", params)
	}

	specificParams, err := ps.GetParameter(ctx, []string{"powGetSysLoad", "cmsBattSoc", "backupReverseSoc"})
	if err != nil {
		slog.Error("Failed to get parameters", "error", err)
	} else {
		slog.Info("Power Station Specific Parameters", "params", specificParams)
	}

	//set functions
	// Example: Set Backup Reserve Level to 24% if it's not already set
	// JSON parser returns numbers as float64 by default so we need to convert
	var currentLevel int
	switch v := specificParams.Data["backupReverseSoc"].(type) {
	case float64:
		currentLevel = int(v)
	case int:
		currentLevel = v
	default:
		slog.Error("Unexpected type for backupReverseSoc", "type", fmt.Sprintf("%T", v))
		return
	}
	slog.Info("Current Backup Reverse SOC Level", "level", currentLevel)

	var backupReverseLevel = 23
	if currentLevel != backupReverseLevel {
		ps.SetBackupReserveLevel(ctx, backupReverseLevel)
	}

	//History data retrieval examples
	beginTime, endTime := time.Now().AddDate(0, 0, -1), time.Now()
	historyData, err := ps.GetBatteryChargingDischargingPower(ctx, beginTime, endTime)
	if err != nil {
		slog.Error("Failed to get history data", "error", err)
	} else {
		slog.Info("Power Station History Data", "data", historyData)
	}

}

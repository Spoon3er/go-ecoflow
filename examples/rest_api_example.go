//go:build exclude

package main

import (
	"context"
	"log/slog"
	"os"

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
	ps := client.GetStream("SN_HERE")

	//set functions
	// ps.SetPowerSocket(ctx, ecoflow.AC1, true)
	// ps.SetBackupReserveLevel(ctx, 20)
	// ps.SetOperatingMode(ctx, ecoflow.SelfPoweredMode, true)

	// get functions
	params, err := ps.GetAllParameters(ctx)
	if err != nil {
		slog.Error("Failed to get parameters", "error", err)
	} else {
		slog.Info("Power Station All Parameters", "params", params)
	}

	specificParams, err := ps.GetParameter(ctx, []string{"powGetSysLoad", "cmsBattSoc"})
	if err != nil {
		slog.Error("Failed to get parameters", "error", err)
	} else {
		slog.Info("Power Station Specific Parameters", "params", specificParams)
	}

	// get set / get functions for smart meters
	sm := client.GetSmartMeter("SN_HERE")
	params, err = sm.GetAllParameters(ctx)
	if err != nil {
		slog.Error("Failed to get parameters", "error", err)
	} else {
		slog.Info("Smart Meter All Parameters", "params", params)
	}
}

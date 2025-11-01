//go:build exclude

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Spoon3er/go-ecoflow"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	// Create Ecoflow client using AccessKey and SecretKey (same as REST API)
	var accessKey = os.Getenv("ACCESS_KEY")
	var secretKey = os.Getenv("SECRET_KEY")
	var deviceSn = os.Getenv("DEVICE_SN")

	client := ecoflow.NewEcoflowClient(accessKey, secretKey)

	// Initialize MQTT using the same client
	mqttConfig := ecoflow.MqttClientConfiguration{
		OnConnect:            connectHandler,
		OnConnectionLost:     connectLostHandler,
		OnReconnect:          reconnectHandler,
		MaxReconnectInterval: time.Hour, // default is 10 minutes
	}

	err := client.InitializeMqtt(context.Background(), mqttConfig)
	if err != nil {
		log.Fatalf("Unable to initialize MQTT: %+v\n", err)
	}

	mqttClient := client.GetMqttClient()

	// Subscribe to device quota topic to receive real-time parameter updates
	// Topic: /open/{certificateAccount}/{sn}/quota
	err = mqttClient.SubscribeDeviceQuota(deviceSn, quotaMessageHandler)
	if err != nil {
		log.Fatalf("Unable to subscribe to quota: %+v\n", err)
	}

	fmt.Println("Listening for device updates...")
	fmt.Println("Press Ctrl+C to exit...")

	// Example: Send a command via MQTT (optional)
	// This demonstrates the set/set_reply pattern
	// Commented out for now to focus on receiving messages
	// go sendExampleCommand(mqttClient, deviceSn)

	// Setup signal handling to gracefully shutdown on Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Block until interrupt signal is received
	<-sigChan
	fmt.Println("\nShutting down gracefully...")

	mqttClient.Disconnect(250)
	fmt.Println("Disconnected from MQTT broker")
}

func sendExampleCommand(client *ecoflow.MqttClient, deviceSn string) {
	// Wait a bit for subscription to be ready
	time.Sleep(5 * time.Second)

	// Example: Send a command with reply
	request := &ecoflow.MqttSetRequest{
		Id:          fmt.Sprintf("%d", time.Now().UnixMilli()),
		Version:     "1.0",
		ModuleType:  1,
		OperateType: "TCP",
		Params: map[string]any{
			"id":      123,
			"enabled": 1,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reply, err := client.PublishSetCommandWithReply(ctx, deviceSn, request, 10*time.Second)
	if err != nil {
		log.Printf("Command failed: %+v\n", err)
		return
	}

	if reply.Code == "0" {
		log.Printf("Command successful: %+v\n", reply)
	} else {
		log.Printf("Command failed with code %s: %s\n", reply.Code, reply.Message)
	}
}

// quotaMessageHandler handles device parameter updates from the quota topic
var quotaMessageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	payload := msg.Payload()

	// Count parameters for summary
	var params ecoflow.MqttDeviceParams
	if err := json.Unmarshal(payload, &params); err != nil {
		fmt.Printf("Unable to parse message from topic %s due to error: %+v\n", msg.Topic(), err)
		return
	}

	// Print timestamp and parameter count
	fmt.Printf("\n[%s] Received %d parameters\n", time.Now().Format("15:04:05"), len(params))

	// Extract power parameters using Stream device types
	powerParams := ecoflow.ExtractStreamPowerParams(payload)
	// Now you can access: powerParams.PowGetSysGrid, etc.

	// Extract BMS parameters
	bmsParams := ecoflow.ExtractStreamBMSParams(payload)
	// Now you can access: bmsParams.BmsBattSoc, etc.

	// Extract CMS parameters
	cmsParams := ecoflow.ExtractStreamCMSParams(payload)
	// Now you can access: cmsParams.CmsBattSoc, etc.

	// Extract PV parameters
	pvParams := ecoflow.ExtractStreamPVParams(payload)
	// Now you can access: pvParams.PlugInInfoPvVol, etc.

	// Extract system parameters
	sysParams := ecoflow.ExtractStreamSystemParams(payload)
	// Now you can access: sysParams.ModuleWifiRssi, etc.

	// Pretty print for display (optional)
	if jsonData, _ := json.Marshal(powerParams); string(jsonData) != "{}" {
		fmt.Printf("PowerParams: %s\n", jsonData)
	}
	if jsonData, _ := json.Marshal(bmsParams); string(jsonData) != "{}" {
		fmt.Printf("BMSParams: %s\n", jsonData)
	}
	if jsonData, _ := json.Marshal(cmsParams); string(jsonData) != "{}" {
		fmt.Printf("CMSParams: %s\n", jsonData)
	}
	if jsonData, _ := json.Marshal(pvParams); string(jsonData) != "{}" {
		fmt.Printf("PVParams: %s\n", jsonData)
	}
	if jsonData, _ := json.Marshal(sysParams); string(jsonData) != "{}" {
		fmt.Printf("SystemParams: %s\n", jsonData)
	}

	// Example: Access the actual typed values
	if powerParams.PowGetSysGrid != nil {
		fmt.Printf("System Grid Power: %.2f W\n", *powerParams.PowGetSysGrid)
	}
	if bmsParams.BmsBattSoc != nil {
		fmt.Printf("Battery State of Charge: %.2f %%\n", *bmsParams.BmsBattSoc)
	}
	if sysParams.ModuleWifiRssi != nil {
		fmt.Printf("WiFi RSSI: %.0f dBm\n", *sysParams.ModuleWifiRssi)
	}
}

// connectHandler executes when successfully connected to the MQTT broker
var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	optionsReader := client.OptionsReader()
	fmt.Printf("Connected to MQTT broker: %s\n", optionsReader.Servers()[0].String())
}

// connectLostHandler executes when the MQTT connection is lost
var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("MQTT connection lost: %v\n", err)
}

// reconnectHandler executes when attempting to reconnect to the MQTT broker
var reconnectHandler mqtt.ReconnectHandler = func(client mqtt.Client, opts *mqtt.ClientOptions) {
	fmt.Println("Attempting to reconnect to MQTT broker...")
}

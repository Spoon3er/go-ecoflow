package ecoflow

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	mqttCertificationUrl = "/iot-open/sign/certification"
)

// MqttClientConfiguration holds configuration for MQTT client initialization
type MqttClientConfiguration struct {
	OnConnect            mqtt.OnConnectHandler
	OnConnectionLost     mqtt.ConnectionLostHandler
	OnReconnect          mqtt.ReconnectHandler
	MaxReconnectInterval time.Duration
}

// MqttClient wraps the MQTT client and provides Ecoflow-specific functionality
type MqttClient struct {
	Client           mqtt.Client
	connectionConfig *MqttConnectionConfig
	pendingReplies   map[int]chan *MqttSetReply // for request/response correlation
	repliesMutex     sync.RWMutex
}

// newMqttClient creates a new MQTT client using the provided connection configuration
// This is an internal function called by Client.InitializeMqtt()
func newMqttClient(connectionConfig *MqttConnectionConfig, config MqttClientConfiguration) *MqttClient {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("%s://%s:%s", connectionConfig.Protocol, connectionConfig.Url, connectionConfig.Port))
	opts.SetClientID(fmt.Sprintf("%s_go-ecoflow", connectionConfig.CertificateAccount))
	opts.SetUsername(connectionConfig.CertificateAccount)
	opts.SetPassword(connectionConfig.CertificatePassword)
	opts.SetConnectRetry(true)

	if config.OnConnect != nil {
		opts.OnConnect = config.OnConnect
	}
	if config.OnConnectionLost != nil {
		opts.OnConnectionLost = config.OnConnectionLost
	}
	if config.OnReconnect != nil {
		opts.OnReconnecting = config.OnReconnect
	}
	// Default value is 10 minutes
	if config.MaxReconnectInterval != 0 {
		opts.MaxReconnectInterval = config.MaxReconnectInterval
	}

	return &MqttClient{
		Client:           mqtt.NewClient(opts),
		connectionConfig: connectionConfig,
		pendingReplies:   make(map[int]chan *MqttSetReply),
	}
}

// getMqttCredentials fetches MQTT credentials using AccessKey/SecretKey authentication
// Endpoint: GET /iot-open/sign/certification
// Uses the same authentication mechanism as HTTP API (accessKey, timestamp, nonce, sign headers)
func getMqttCredentials(ctx context.Context, client *Client) (*MqttConnectionConfig, error) {
	httpReq := NewHttpRequest(
		client.httpClient,
		"GET",
		client.baseUrl+mqttCertificationUrl,
		nil,
		client.accessToken,
		client.secretToken,
	)

	responseBody, err := httpReq.Execute(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get MQTT credentials: %w", err)
	}

	var mqttCreds MqttCredentialsResponse
	err = json.Unmarshal(responseBody, &mqttCreds)
	if err != nil {
		return nil, fmt.Errorf("failed to parse MQTT credentials response: %w", err)
	}

	if mqttCreds.Code != "0" {
		return nil, fmt.Errorf("failed to get MQTT credentials: code=%s, message=%s", mqttCreds.Code, mqttCreds.Message)
	}

	return &mqttCreds.Data, nil
}

// Connect connects to the MQTT broker
func (m *MqttClient) Connect() error {
	if token := m.Client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// Disconnect disconnects from the MQTT broker
func (m *MqttClient) Disconnect(quiesce uint) {
	m.Client.Disconnect(quiesce)
}

// SubscribeDeviceQuota subscribes to the device quota topic to receive real-time parameter updates
// Topic: /open/{certificateAccount}/{deviceSn}/quota
// The client must be connected to the broker before subscribing
func (m *MqttClient) SubscribeDeviceQuota(deviceSn string, callback mqtt.MessageHandler) error {
	topic := fmt.Sprintf("/open/%s/%s/quota", m.connectionConfig.CertificateAccount, deviceSn)
	token := m.Client.Subscribe(topic, 1, callback)
	token.Wait()
	if token.Error() != nil {
		return fmt.Errorf("failed to subscribe to quota topic: %w", token.Error())
	}
	return nil
}

// SubscribeSetReply subscribes to the set_reply topic to receive command responses
// Topic: /open/{certificateAccount}/{deviceSn}/set_reply
// This is automatically called when using PublishSetCommandWithReply
func (m *MqttClient) SubscribeSetReply(deviceSn string) error {
	topic := fmt.Sprintf("/open/%s/%s/set_reply", m.connectionConfig.CertificateAccount, deviceSn)

	handler := func(client mqtt.Client, msg mqtt.Message) {
		var reply MqttSetReply
		if err := json.Unmarshal(msg.Payload(), &reply); err != nil {
			return
		}

		m.repliesMutex.RLock()
		ch, exists := m.pendingReplies[reply.Id]
		m.repliesMutex.RUnlock()

		if exists {
			select {
			case ch <- &reply:
			default:
			}
		}
	}

	token := m.Client.Subscribe(topic, 1, handler)
	token.Wait()
	if token.Error() != nil {
		return fmt.Errorf("failed to subscribe to set_reply topic: %w", token.Error())
	}
	return nil
}

// PublishSetCommand publishes a command to the device (fire and forget, no reply)
// Topic: /open/{certificateAccount}/{deviceSn}/set
func (m *MqttClient) PublishSetCommand(deviceSn string, request map[string]any) error {
	topic := fmt.Sprintf("/open/%s/%s/set", m.connectionConfig.CertificateAccount, deviceSn)

	payload, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	token := m.Client.Publish(topic, 1, false, payload)
	token.Wait()
	if token.Error() != nil {
		return fmt.Errorf("failed to publish command: %w", token.Error())
	}
	return nil
}

// PublishSetCommandWithReply publishes a command and waits for a reply
// Topic: /open/{certificateAccount}/{deviceSn}/set
// Reply on: /open/{certificateAccount}/{deviceSn}/set_reply
func (m *MqttClient) PublishSetCommandWithReply(ctx context.Context, deviceSn string, request map[string]any, timeout time.Duration) (*MqttSetReply, error) {
	// Ensure we're subscribed to set_reply topic
	if err := m.SubscribeSetReply(deviceSn); err != nil {
		return nil, err
	}

	// Create reply channel
	replyCh := make(chan *MqttSetReply, 1)

	// Safe type assertion for request ID
	var requestId int
	if idValue, exists := request["id"]; exists && idValue != nil {
		if id, ok := idValue.(int); ok {
			requestId = id
		} else {
			return nil, fmt.Errorf("request 'id' field must be an int, got %T", idValue)
		}
	} else {
		return nil, fmt.Errorf("request 'id' field is required and cannot be nil")
	}

	m.repliesMutex.Lock()
	m.pendingReplies[requestId] = replyCh
	m.repliesMutex.Unlock()

	// Cleanup
	defer func() {

		m.repliesMutex.Lock()
		delete(m.pendingReplies, requestId)
		m.repliesMutex.Unlock()
		close(replyCh)
	}()

	// Publish command
	if err := m.PublishSetCommand(deviceSn, request); err != nil {
		return nil, err
	}

	// Wait for reply or timeout
	select {
	case reply := <-replyCh:
		return reply, nil
	case <-time.After(timeout):
		return nil, fmt.Errorf("timeout waiting for command reply")
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

package ecoflow

// MqttConnectionConfig represents the configuration for MQTT connection
// returned from /iot-open/sign/certification endpoint.
// It contains the following fields:
// - CertificateAccount: the username for MQTT authentication (also used in topic paths)
// - CertificatePassword: the password for MQTT authentication
// - Url: the URL of the MQTT broker
// - Port: the port number for the MQTT connection
// - Protocol: the protocol for the MQTT connection (e.g., "mqtts")
type MqttConnectionConfig struct {
	CertificateAccount  string `json:"certificateAccount"`
	CertificatePassword string `json:"certificatePassword"`
	Url                 string `json:"url"`
	Port                string `json:"port"`
	Protocol            string `json:"protocol"`
}

// MqttCredentialsResponse represents the response structure for MQTT credentials
// from GET /iot-open/sign/certification endpoint.
// Example response:
//
//	{
//	  "code": "0",
//	  "message": "Success",
//	  "data": {
//	    "certificateAccount": "open-57c134518b5***",
//	    "certificatePassword": "959253cc103a4008***",
//	    "url": "mqtt.ecoflow.com",
//	    "port": "8883",
//	    "protocol": "mqtts"
//	  }
//	}
type MqttCredentialsResponse struct {
	Code    string               `json:"code"`
	Message string               `json:"message"`
	Data    MqttConnectionConfig `json:"data"`
}

// MqttDeviceParams represents the device parameters received from MQTT quota topic
// Topic: /open/{certificateAccount}/{sn}/quota
// The payload is a flat JSON object with parameter names as keys
// Example:
// {"powGetSysGrid":409.0,"powGetSysLoad":409.0,"bmsDsgRemTime":4211,"bmsMaxCellTemp":20,...}
// Or for more complex nested data:
// {"soc":22,"num":0,"cellVol":[3259,3258,3257],"cellTemp":[20,19],...}
type MqttDeviceParams map[string]any

// MqttSetReply represents the response received from MQTT set_reply topic
// Topic: /open/{certificateAccount}/{sn}/set_reply
// The ID field matches the request ID for correlation
// Example:
//
//	{
//	  "id": 1681872798000,
//	  "code": "0",
//	  "message": "Success"
//	}
type MqttSetReply struct {
	Id      int    `json:"id"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

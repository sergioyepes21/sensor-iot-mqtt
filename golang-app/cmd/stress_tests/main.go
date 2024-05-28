package stresstests

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	brokerpublisher "github.com/sergioyepes21/sensor-iot-mqtt/internal/broker-publisher"
	custombroker "github.com/sergioyepes21/sensor-iot-mqtt/internal/custom-broker"
)

var mqttPublisher = brokerpublisher.NewMQTTPublisher()
var mqttBroker = custombroker.NewMQTTBroker(mqttPublisher, nil)

func main() {
	// Your code here
	mqttPublisher.Publish(mqtt.NewClient())
}

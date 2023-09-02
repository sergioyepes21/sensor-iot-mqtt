package brokerconsumer

import (
	"encoding/binary"
	"fmt"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	anomalynotification "github.com/sergioyepes21/sensor-iot-mqtt/cmd/sensor-anomaly/anomaly-notification"
)

func bytesToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	return float64(bits)
}

type MQTTConsumer struct {
	anomalyNotification *anomalynotification.AnomalyNotification
}

func NewMQTTConsumer(a *anomalynotification.AnomalyNotification) *MQTTConsumer {
	return &MQTTConsumer{
		anomalyNotification: a,
	}
}

func (c *MQTTConsumer) Consume(client mqtt.Client, msg mqtt.Message, wg *sync.WaitGroup) {
	sensorData, err := c.fromMQTTMessageToSensorData(msg)
	if err != nil {
		fmt.Printf("Error processing Message %d: %v\n", msg.MessageID(), err)
		return
	}

	anomalousData := c.getAnomalousSensors(*sensorData)
	anomalyDetected := len(*anomalousData) > 0

	if anomalyDetected {
		go c.anomalyNotification.Notify(sensorData.VehicleId, sensorData.Latitude, sensorData.Longitude, anomalousData)
	}

	wg.Done()
}

func (c *MQTTConsumer) fromMQTTMessageToSensorData(msg mqtt.Message) (*SensorData, error) {
	rawBytes := msg.Payload()

	sensorData := SensorData{
		VehicleId: string(rawBytes[0:36]),
		Latitude:  bytesToFloat64(rawBytes[36:52]),
		Longitude: bytesToFloat64(rawBytes[52:68]),
		SensoredValues: map[string]bool{
			"Speed":       rawBytes[68] > 0,
			"Temperature": rawBytes[69] > 0,
			"Humidity":    rawBytes[70] > 0,
			"Stops":       rawBytes[71] > 0,
			"Brakes":      rawBytes[72] > 0,
		},
	}

	return &sensorData, nil
}

func (c *MQTTConsumer) getAnomalousSensors(sensorData SensorData) *[]string {
	anomalyOnSensor := new([]string)
	for sensorKey, sensorValue := range sensorData.SensoredValues {
		if sensorValue {
			*anomalyOnSensor = append(*anomalyOnSensor, sensorKey)
		}
	}
	return anomalyOnSensor
}

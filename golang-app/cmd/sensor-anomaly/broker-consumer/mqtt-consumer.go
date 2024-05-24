package brokerconsumer

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	anomalynotification "github.com/sergioyepes21/sensor-iot-mqtt/cmd/sensor-anomaly/anomaly-notification"
	logger "github.com/sergioyepes21/sensor-iot-mqtt/internal/custom-logger"
)

type MQTTConsumer struct {
	anomalyNotification *anomalynotification.AnomalyNotification
}

func NewMQTTConsumer(a *anomalynotification.AnomalyNotification) *MQTTConsumer {
	return &MQTTConsumer{
		anomalyNotification: a,
	}
}

func (c *MQTTConsumer) Consume(client mqtt.Client, msg mqtt.Message, wg *sync.WaitGroup) {

	startTime := time.Now()

	sensorData, err := c.fromMQTTMessageToSensorData(msg)
	if err != nil {
		fmt.Printf("Error processing Message %d: %v\n", msg.MessageID(), err)
		return
	}

	anomalyDetected := c.getAnomalousSensors(*sensorData)

	if anomalyDetected {
		go c.anomalyNotification.Notify(sensorData.VehicleId, sensorData.Latitude, sensorData.Longitude, startTime)
	} else {
		endTime := time.Now()
		duration := endTime.Sub(startTime)
		logger.Log.Printf("No anomalies on [Vehicle: %s], [Duration: %v]", sensorData.VehicleId, duration.Microseconds())
	}

	wg.Done()
}

func (c *MQTTConsumer) fromMQTTMessageToSensorData(msg mqtt.Message) (*SensorData, error) {
	rawBytes := msg.Payload()

	messageStr := string(rawBytes[:])
	messageArrValues := strings.Split(messageStr, " ")

	if len(messageArrValues) != 8 {
		return nil, fmt.Errorf("message %s has not the expected format", messageStr)
	}

	latitude, err := getLatitudeFromBrokerMessage(messageArrValues)
	if err != nil {
		return nil, err
	}

	longitude, err := getLongitudeFromBrokerMessage(messageArrValues)
	if err != nil {
		return nil, err
	}

	sensorData := SensorData{
		VehicleId: messageArrValues[0],
		Latitude:  latitude,
		Longitude: longitude,
		SensoredValues: map[string]bool{
			"Speed":       messageArrValues[3] == "1",
			"Temperature": messageArrValues[4] == "1",
			"Humidity":    messageArrValues[5] == "1",
			"Stops":       messageArrValues[6] == "1",
			"Brakes":      messageArrValues[7] == "1",
		},
	}

	return &sensorData, nil
}

func getLatitudeFromBrokerMessage(messageArr []string) (float64, error) {
	latStr := messageArr[1]
	latitude, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing latitude %s: %v", latStr, err)
	}
	return latitude, nil
}

func getLongitudeFromBrokerMessage(messageArr []string) (float64, error) {
	latStr := messageArr[2]
	latitude, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing longitude %s: %v", latStr, err)
	}
	return latitude, nil
}

func (c *MQTTConsumer) getAnomalousSensors(sensorData SensorData) bool {
	for _, sensorValue := range sensorData.SensoredValues {
		if sensorValue {
			return true
		}
	}
	return false
}

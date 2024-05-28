package brokerpublisher

import (
	"fmt"
	"strings"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"golang.org/x/exp/rand"
)

type MQTTPublisher struct {
}

func NewMQTTPublisher() *MQTTPublisher {
	return &MQTTPublisher{}
}

func (p *MQTTPublisher) Publish(c mqtt.Client, wg *sync.WaitGroup) {
	defer wg.Done()

	sensorData := createSensorDataStr()

	topic := createRandomTopicId()

	token := c.Publish(topic, 0, false, sensorData)
	token.Wait()
}

func createSensorDataStr() string {
	messageArr := []string{"e6f42c2e-4733-11ee-be56-0242ac120002", "37.7749", "-122.4194", "1", "0", "1", "0", "1"}
	return strings.Join(messageArr, " ")
}

func crateRandomVehicleId() string {

}

func createRandomTopicId() string {
	randomInt := rand.Intn(100)
	return fmt.Sprintf("testtopic/%d", randomInt)
}

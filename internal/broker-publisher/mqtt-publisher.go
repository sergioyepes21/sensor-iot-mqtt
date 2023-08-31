package brokerpublisher

import (
	"encoding/binary"
	"math"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTPublisher struct {
}

func NewMQTTPublisher() *MQTTPublisher {
	return &MQTTPublisher{}
}

func (p *MQTTPublisher) Publish(c mqtt.Client, wg *sync.WaitGroup) {
	defer wg.Done()

	sensorData := createSensorDataBytesArray()

	token := c.Publish("testtopic/1", 0, false, sensorData)
	token.Wait()
}

func createSensorDataBytesArray() []byte {
	byteData := make([]byte, 73)
	copy(byteData[0:36], []byte("e6f42c2e-4733-11ee-be56-0242ac120002"))
	copyFloat64ToBytes(byteData[36:52], 37.7749)
	copyFloat64ToBytes(byteData[52:68], -122.4194)
	byteData[68] = 1
	byteData[69] = 0
	byteData[70] = 1
	byteData[71] = 0
	byteData[72] = 1
	return byteData
}

func copyFloat64ToBytes(dest []byte, value float64) {
	bits := math.Float64bits(value)
	binary.LittleEndian.PutUint64(dest, bits)
}

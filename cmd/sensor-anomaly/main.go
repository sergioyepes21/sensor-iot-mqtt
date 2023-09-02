package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	brokerpublisher "github.com/sergioyepes21/sensor-iot-mqtt/internal/broker-publisher"
	custombroker "github.com/sergioyepes21/sensor-iot-mqtt/internal/custom-broker"

	anomalynotification "github.com/sergioyepes21/sensor-iot-mqtt/cmd/sensor-anomaly/anomaly-notification"
	brokerconsumer "github.com/sergioyepes21/sensor-iot-mqtt/cmd/sensor-anomaly/broker-consumer"
)

var consumerWG = new(sync.WaitGroup)

// var publisherWG = new(sync.WaitGroup)

var anomalyNotification = anomalynotification.NewAnomalyNotification()

var mqttPublisher = brokerpublisher.NewMQTTPublisher()
var mqttConsumer = brokerconsumer.NewMQTTConsumer(anomalyNotification)
var mqttBroker = custombroker.NewMQTTBroker(mqttPublisher, mqttConsumer)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	consumerWG.Add(1)
	go mqttConsumer.Consume(client, msg, consumerWG)
}

func main() {

	// c, err := mqttBroker.Start(f)
	_, err := mqttBroker.Start(f)
	if err != nil {
		panic(err)
	}

	// mqttBroker.PublishTestData(c, publisherWG)
	// publisherWG.Wait()

	timeToWait := 1 * time.Second
	time.Sleep(timeToWait)
	waitForExitSignal()
}

func waitForExitSignal() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan
	fmt.Println("Exiting...")
}

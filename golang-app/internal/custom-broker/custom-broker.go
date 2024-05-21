package custombroker

import (
	"fmt"
	"os"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

type MQTTBrokerPublisherInterface interface {
	Publish(c mqtt.Client, wg *sync.WaitGroup)
}

type MQTTBrokerConsumerInterface interface {
	Consume(client mqtt.Client, msg mqtt.Message, wg *sync.WaitGroup)
}

type MQTTBroker struct {
	mqttPublisher  MQTTBrokerPublisherInterface
	mqttConsumer   MQTTBrokerConsumerInterface
	brokerHost     string
	brokerClientId string
	brokerTopic    string
}

func NewMQTTBroker(mp MQTTBrokerPublisherInterface, mc MQTTBrokerConsumerInterface) *MQTTBroker {
	cliendIdValue := uuid.New()
	cliendIdString := cliendIdValue.String()

	brokerHost := getEnv("BROKER_HOST", "tcp://mqtt5:1883")
	brokerClientId := getEnv("BROKER_CLIENT_ID", cliendIdString)
	brokerTopic := getEnv("BROKER_TOPIC", "testtopic/#")

	fmt.Printf("Broker host: %s ~ Broker client id: %s ~ Broker topic: %s\n", brokerHost, brokerClientId, brokerTopic)

	return &MQTTBroker{
		mqttPublisher:  mp,
		mqttConsumer:   mc,
		brokerHost:     brokerHost,
		brokerClientId: brokerClientId,
		brokerTopic:    brokerTopic,
	}
}

func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func (b *MQTTBroker) Start(f mqtt.MessageHandler) (mqtt.Client, error) {
	opts := b.createMQTTClientOptions(f)
	c := mqtt.NewClient(opts)

	err := b.connectToTopic(c)
	if err != nil {
		return nil, err
	}

	err = b.subscribeToTopic(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (b *MQTTBroker) createMQTTClientOptions(f mqtt.MessageHandler) *mqtt.ClientOptions {
	// mqtt.DEBUG = log.New(os.Stdout, "", 0)
	// mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker(b.brokerHost).SetClientID(b.brokerClientId).SetUsername("user1").SetPassword("user1")

	opts.SetKeepAlive(60 * time.Second)
	// Set the message callback handler
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)
	opts.WillQos = 1

	return opts
}

func (b *MQTTBroker) connectToTopic(c mqtt.Client) error {
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (b *MQTTBroker) subscribeToTopic(c mqtt.Client) error {
	if token := c.Subscribe(b.brokerTopic, 0, nil); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (b *MQTTBroker) PublishTestData(c mqtt.Client, wg *sync.WaitGroup) {
	numIterations := 1
	// Publish a message
	for i := 0; i < numIterations; i++ {
		wg.Add(1)
		go b.mqttPublisher.Publish(c, wg)
	}
	wg.Wait()
}

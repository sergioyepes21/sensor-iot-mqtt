package custombroker

import (
	"log"
	"os"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const BrokerHost = "tcp://127.0.01:1883"
const BrokerClientId = "emqx_test_client"
const BrokerTopic = "testtopic/#"

type MQTTBrokerPublisherInterface interface {
	Publish(c mqtt.Client, wg *sync.WaitGroup)
}

type MQTTBrokerConsumerInterface interface {
	Consume(client mqtt.Client, msg mqtt.Message, wg *sync.WaitGroup)
}

type MQTTBroker struct {
	mqttPublisher MQTTBrokerPublisherInterface
	mqttConsumer  MQTTBrokerConsumerInterface
}

func NewMQTTBroker(mp MQTTBrokerPublisherInterface, mc MQTTBrokerConsumerInterface) *MQTTBroker {
	return &MQTTBroker{
		mqttPublisher: mp,
		mqttConsumer:  mc,
	}
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
	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker(BrokerHost).SetClientID(BrokerClientId)

	opts.SetKeepAlive(60 * time.Second)
	// Set the message callback handler
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)
	return opts
}

func (b *MQTTBroker) connectToTopic(c mqtt.Client) error {
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
func (b *MQTTBroker) subscribeToTopic(c mqtt.Client) error {
	if token := c.Subscribe(BrokerTopic, 0, nil); token.Wait() && token.Error() != nil {
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

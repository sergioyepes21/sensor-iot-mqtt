package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

func main() {
	opts := createMQTTClientOptions()
	c := mqtt.NewClient(opts)

	err := connectToTopic(c)
	if err != nil {
		panic(err)
	}

	publishMessages(c)
}

func createMQTTClientOptions() *mqtt.ClientOptions {
	// brokerHost := getEnv("BROKER_HOST", "tcp://mqtt5:1883")
	brokerHost := getEnv("BROKER_HOST", "tcp://localhost:1883")
	brokerClientId := getEnv("BROKER_CLIENT_ID", uuid.NewString())
	// brokerTopic := getEnv("BROKER_TOPIC", "testtopic/#")

	username := getEnv("BROKER_USERNAME", "user1")
	password := getEnv("BROKER_PASSWORD", "user1")

	opts := mqtt.NewClientOptions().AddBroker(brokerHost).SetClientID(brokerClientId).SetUsername(username).SetPassword(password)

	opts.SetKeepAlive(60 * time.Second)

	opts.SetPingTimeout(1 * time.Second)
	opts.WillQos = 1

	return opts
}

func getEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func connectToTopic(c mqtt.Client) error {
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func publishMessages(c mqtt.Client) {

	startTime := time.Now()

	testDurationInMinutes := 15
	messagesPerMinute := 60

	var iteration int

	for time.Since(startTime) < time.Duration(testDurationInMinutes)*time.Minute {
		publishMessagesPerMinute(c, messagesPerMinute)
		iteration++
		fmt.Printf("Iteration %d finished\n out of %d total", iteration, testDurationInMinutes)
	}

	fmt.Println("Test finished")
}

func publishMessagesPerMinute(c mqtt.Client, messagesPerMinute int) {
	var publisherWG = new(sync.WaitGroup)
	waitPeriod := 60 / messagesPerMinute
	for i := 0; i < messagesPerMinute; i++ {
		publisherWG.Add(1)
		go publishMessage(c, publisherWG)
		time.Sleep(time.Duration(waitPeriod) * time.Second)
	}

	publisherWG.Wait()
}

func publishMessage(c mqtt.Client, wg *sync.WaitGroup) {
	message := createSensorDataStr()
	topic := createRandomTopicId()

	token := c.Publish(topic, 0, false, message)
	token.Wait()
	defer wg.Done()
}

func createSensorDataStr() string {
	vehicleId := createRandomVehicleId()
	latitude := rand.Float64()
	longitude := rand.Float64()

	// Total of probability of anomalies is 0.5
	speed := isAnomalyWithProbability(0.1)
	temperature := isAnomalyWithProbability(0.1)
	humidity := isAnomalyWithProbability(0.1)
	stops := isAnomalyWithProbability(0.1)
	breaks := isAnomalyWithProbability(0.1)

	return fmt.Sprintf("%s %f %f %d %d %d %d %d", vehicleId, latitude, longitude, speed, temperature, humidity, stops, breaks)
}

func createRandomVehicleId() string {
	fixedVehicleIds := []string{"e6f42c2e-4733-11ee-be56-0242ac120002", "b85f4f57-9e70-4b93-9d9f-be25467dbbd7", "bdbf886e-b5e6-4c97-97eb-031681f6cc7b"}

	randomInt := rand.Float64()

	vehicleId := uuid.NewString()

	if randomInt < 0.7 {
		randomInt := rand.Intn(len(fixedVehicleIds) - 1)
		vehicleId = fixedVehicleIds[randomInt]
	}

	return vehicleId
}

func isAnomalyWithProbability(probability float64) int {
	randomFloat := rand.Float64()
	isAnomaly := randomFloat < probability

	if isAnomaly {
		return 1
	}
	return 0
}

func createRandomTopicId() string {
	randomInt := rand.Intn(100)
	return fmt.Sprintf("testtopic/%d", randomInt)
}

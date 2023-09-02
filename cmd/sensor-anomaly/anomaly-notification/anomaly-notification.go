package anomalynotification

import (
	"fmt"
	"time"

	logger "github.com/sergioyepes21/sensor-iot-mqtt/internal/custom-logger"
	redisclient "github.com/sergioyepes21/sensor-iot-mqtt/internal/redis-client"
)

type AnomalyNotification struct {
	redisClient *redisclient.MyRedisClient
}

func NewAnomalyNotification() *AnomalyNotification {
	redisClient := redisclient.NewMyRedisClient()

	return &AnomalyNotification{
		redisClient: redisClient,
	}
}

func (a *AnomalyNotification) Notify(vehicleId string, lat float64, long float64, anomalousData *[]string, startTime time.Time) {
	values, err := a.redisClient.GetHashValues(vehicleId)
	if err != nil {
		fmt.Printf("Error getting key %s: %v\n", vehicleId, err)
		return
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	logger.Log.Printf("[Duration: %v], [Vehicle: %s], [Phone numbers: %s], [Latitude: %f], [Longitude: %f], [Anomalous data: %v]\n", duration, vehicleId, values, lat, long, anomalousData)
}

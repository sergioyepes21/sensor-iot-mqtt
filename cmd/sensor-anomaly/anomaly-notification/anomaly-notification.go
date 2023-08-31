package anomalynotification

import (
	"fmt"

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

func (a *AnomalyNotification) Notify(vehicleId string, lat float64, long float64, anomalousData *[]string) {
	values, err := a.redisClient.GetHashValues(vehicleId)
	if err != nil {
		fmt.Printf("Error getting key %s: %v\n", vehicleId, err)
		return
	}

	fmt.Printf("Alert on Vehicle: %s ~ Notifying the following numbers: %s, With latitude %f, longitude%f\n", vehicleId, values, lat, long)
}

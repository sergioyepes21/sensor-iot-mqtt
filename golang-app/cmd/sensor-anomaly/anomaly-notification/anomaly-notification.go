package anomalynotification

import (
	"time"

	logger "github.com/sergioyepes21/sensor-iot-mqtt/internal/custom-logger"
	redisclient "github.com/sergioyepes21/sensor-iot-mqtt/internal/redis-client"
)

type AnomalyNotification struct {
}

func NewAnomalyNotification() *AnomalyNotification {

	return &AnomalyNotification{}
}

func (a *AnomalyNotification) Notify(redisClient *redisclient.MyRedisClient, vehicleId string, lat float64, long float64, startTime time.Time) {
	// redisClient.Conn(redisClient.Context())

	vehicleContact, err := redisClient.GetHashValues(vehicleId)

	// errClose := redisClient.Close()
	if err != nil {
		logger.Log.Printf("Error getting key %s: ~ GetHV E: %v\n", vehicleId, err)
		return
	}

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	logger.Log.Printf("Anomaly detected on [Vehicle: %s], [Duration: %v], [Phone numbers: %s]", vehicleId, duration.Microseconds(), vehicleContact)
}

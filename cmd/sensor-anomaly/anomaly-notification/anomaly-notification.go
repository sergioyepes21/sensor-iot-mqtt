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

func (a *AnomalyNotification) Notify(vehicleId string, lat float64, long float64, anomalousData *[]string, startTime time.Time) {
	redisClient := redisclient.NewMyRedisClient()
	values, err := redisClient.GetHashValues(vehicleId)
	errClose := redisClient.Close()
	if err != nil || errClose != nil {
		logger.Log.Printf("Error getting key %s: ~ GetHV E: %v ~ CloseConn: %v\n", vehicleId, err, errClose)
		return
	}
	redisClient.Close()

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	logger.Log.Printf("[Duration: %v], [Vehicle: %s], [Phone numbers: %s], [Latitude: %f], [Longitude: %f], [Anomalous data: %v]\n", duration, vehicleId, values, lat, long, anomalousData)
}

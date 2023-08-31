package brokerconsumer

type SensorData struct {
	VehicleId      string          `json:"vehicleId"`
	Latitude       float64         `json:"latitude"`
	Longitude      float64         `json:"longitude"`
	SensoredValues map[string]bool `json:"sensoredValues"`
}

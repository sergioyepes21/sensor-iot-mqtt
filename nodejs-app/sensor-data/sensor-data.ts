
export class SensorData {
  vehicleId: string;
  latitude: number;
  longitude: number;
  sensoredValues: Record<string, boolean>;

  constructor(vehicleId: string, latitude: number, longitude: number, sensoredValues: Record<string, boolean>) {
    this.vehicleId = vehicleId;
    this.latitude = latitude;
    this.longitude = longitude;
    this.sensoredValues = sensoredValues;
  }

  anomalyDetected(): boolean {
    const values = Object.values(this.sensoredValues);

    return values.some((value) => value === true);
  }
}
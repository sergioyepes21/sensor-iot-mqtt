import { RedisService } from "../cache-service/redis-service";
import { ILoggerService } from "../logger/logger-service";
import { WinstonLoggerService } from "../logger/winston-logger-service";
import { SensorData } from "../sensor-data/sensor-data";
import type { IMQTTBrokerConsumer } from "./mqtt-broker-consumer-interface";

export class MQTTBrokerConsumer implements IMQTTBrokerConsumer {

  private readonly loggerService: ILoggerService;

  constructor() {
    this.loggerService = new WinstonLoggerService();
  }

  async consume(msg: Buffer): Promise<void> {
    const startTime = performance.now();

    let message: string = "";

    try {
      const sensorData = this.fromMQTTMessageToSensorData(msg);

      const anomalyDetected = sensorData.anomalyDetected();

      message = anomalyDetected ? await this.buildAnomalyMessage(sensorData) : this.buildNoAnomalyMessage(sensorData);

    } catch (error) {
      console.error("error processing message", error);
    } finally {
      const endTime = performance.now();
      const duration = endTime - startTime;
      this.loggerService.log(message.replace("$duration", (duration * 1000).toFixed(0)));
    }
  }

  private fromMQTTMessageToSensorData(msg: Buffer): SensorData {

    const msgStr = msg.toString();

    const arrValues = msgStr.split(" ");

    if (arrValues.length !== 8) {
      throw new Error(`message: ${msgStr} has not the expected format`);
    }

    const vehicleId = arrValues[0];

    return new SensorData(
      vehicleId,
      this.getLatitudFromMessage(arrValues),
      this.getLongitudeFromMessage(arrValues),
      {
        speed: arrValues[3] === "1",
        temperature: arrValues[4] === "1",
        humidity: arrValues[5] === "1",
        stops: arrValues[6] === "1",
        brakes: arrValues[7] === "1",
      }
    );

  }

  private getLatitudFromMessage(arrValues: string[]): number {
    const latitude = parseFloat(arrValues[1]);
    if (isNaN(latitude)) {
      throw new Error(`error parsing latitude: ${arrValues[0]}`);
    }

    return latitude;
  }

  private getLongitudeFromMessage(arrValues: string[]): number {
    const longitude = parseFloat(arrValues[2]);
    if (isNaN(longitude)) {
      throw new Error(`error parsing longitude: ${arrValues[1]}`);
    }

    return longitude;
  }

  private async buildAnomalyMessage(sensorData: SensorData): Promise<string> {
    let message = "";
    const cacheService = new RedisService();
    try {
      await cacheService.connect();

      const hashValues = await cacheService.getHashValues(sensorData.vehicleId);

      const phoneNumbers = hashValues.join(" ");

      message = `Anomaly detected on [Vehicle: ${sensorData.vehicleId}], [Duration: $duration], [Phone numbers: [${phoneNumbers}]]`;
    } catch (err) {
      console.error(`Error getting key: ${sensorData.vehicleId}`, err);
    } finally {
      await cacheService.close();
    }

    return message;
  }

  private buildNoAnomalyMessage(sensorData: SensorData): string {
    return `No anomalies on [Vehicle: ${sensorData.vehicleId}], [Duration: $duration]`;
  }
}
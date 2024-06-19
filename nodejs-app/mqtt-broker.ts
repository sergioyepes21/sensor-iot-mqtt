import { connect, IClientOptions, MqttClient } from "mqtt";
import { MQTTBrokerConsumer } from "./mqtt-broker-consumer/mqtt-broker-consumer";
import { IMQTTBrokerConsumer } from "./mqtt-broker-consumer/mqtt-broker-consumer-interface";
import { randomUUID } from "crypto";

export class MQTTBroker {

  private mqttClient: MqttClient | undefined;

  private readonly mqttBrokerConsumer: IMQTTBrokerConsumer;


  constructor() {
    this.mqttBrokerConsumer = new MQTTBrokerConsumer();

  }

  start(): void {
    const client = this.getMqttClient();

    client.on("connect", async () => {

      client.subscribe("testtopic/#", (err: unknown) => {
        if (err) {
          console.error(err);
          throw new Error("Failed to subscribe to topic");
        }
      });

      client.on("message", (_topic: string, message: Buffer) => {
        this.mqttBrokerConsumer.consume(message);
      });
    });

    client.on("error", (err: Error) => {
      console.error(err);
      throw new Error("Failed to connect to MQTT broker");
    })
  }

  private getMqttClient(): MqttClient {
    if (this.mqttClient) return this.mqttClient;

    const options: IClientOptions = {
      host: process.env.MQTT_HOST,
      port: 1883,
      username: process.env.MQTT_USERNAME,
      password: process.env.MQTT_PASSWORD,
      clientId: randomUUID(),
      keepalive: 30000,
      reconnectPeriod: 0,
    };
    this.mqttClient = connect(options);

    return this.mqttClient;
  }

}
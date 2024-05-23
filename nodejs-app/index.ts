import "dotenv/config";
import { MQTTBroker } from "./mqtt-broker";

const mqttBroker = new MQTTBroker();

void mqttBroker.start();
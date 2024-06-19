export interface IMQTTBrokerConsumer {
  consume(msg: Buffer): Promise<void>
}
import { Redis, RedisOptions } from "ioredis";
import { ICacheService } from "./cache-service.interface";

export class RedisService implements ICacheService {

  private redisClient: Redis;

  constructor() {
    const options: RedisOptions = {
      host: process.env.REDIS_HOST,
      lazyConnect: true,
    };

    this.redisClient = new Redis(options);
  }

  async connect(): Promise<void> {
    await this.redisClient.connect();
  }

  async getHashValues(key: string): Promise<string[]> {
    const values = await this.redisClient.hvals(key);
    return values;
  }

  async close(): Promise<void> {
    this.redisClient.disconnect();
  }
}
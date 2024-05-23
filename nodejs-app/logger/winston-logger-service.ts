import winston from "winston";
import { ILoggerService } from "./logger-service";
import path from "path";

export class WinstonLoggerService implements ILoggerService {

  private readonly logger: winston.Logger;

  constructor() {
    const filename = path.join(__dirname, "..", "assets", "node-anomalies.log")
    this.logger = winston.createLogger({
      level: "info",
      levels: winston.config.npm.levels,
      format: winston.format.combine(
        winston.format.timestamp({ format: "YYYY/MM/DD HH:mm:ss" }),
        winston.format.printf(log => `${log.timestamp} ${log.message}`),
      ),
      transports: [
        new winston.transports.File({
          level: "info",
          filename: filename,
          options: {
            flags: "w"
          }
        }),
      ]
    });
  }

  log(message: string): void {
    this.logger.info(message);
  }
}
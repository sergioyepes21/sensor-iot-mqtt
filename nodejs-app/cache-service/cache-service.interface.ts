
export interface ICacheService {

  connect(): Promise<void>;
  
  getHashValues(key: string): Promise<string[]>;

  close(): Promise<void>;
}
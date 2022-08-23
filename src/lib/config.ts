import { get, post } from "./fetch";

export interface HTTPSConfig {
  CertFile: string;
  KeyFile: string;
}
export interface ServiceConfig {
  Target: string;
  Path: string;
}

export interface DynamicServiceConfig {
  Path: string;
  Query: string;
}

export interface Config {
  Service: ServiceConfig[];
  DynamicService: DynamicServiceConfig;
  HTTPS: HTTPSConfig;
}

export const readConfig: () => Promise<Config | null> = async () => {
  let res = await get("/api/config");
  try {
    return JSON.parse(res);
  } catch {
    return null;
  }
};

export const writeConfig: (config: Config) => Promise<void> = async (
  config
) => {
  await post("/api/config", config);
};

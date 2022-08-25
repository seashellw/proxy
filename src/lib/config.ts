import { useDebounceFn } from "@vueuse/core";
import { get, post } from "./fetch";

export interface HTTPSConfig {
  CertFile?: string;
  KeyFile?: string;
}
export interface ServiceConfig {
  Target?: string;
  Path?: string;
}

export interface DynamicServiceConfig {
  Path?: string;
  Query?: string;
}

export interface FileServiceConfig {
  Path?: string;
  Dir?: string;
}

export interface Config {
  Service?: ServiceConfig[];
  FileService?: FileServiceConfig[];
  DynamicService?: DynamicServiceConfig;
  HTTPS?: HTTPSConfig;
}

export const readConfig: () => Promise<Config | {}> = async () => {
  let res = await get("/api/config");
  try {
    return JSON.parse(res);
  } catch {
    return {};
  }
};

export const writeConfig = useDebounceFn(async (config: Config) => {
  await post("/api/config", config);
}, 1000);

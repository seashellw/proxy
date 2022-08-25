import { useDebounceFn } from "@vueuse/core";
import { defineStore } from "pinia";
import { ref } from "vue";
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

export const useConfigStore = defineStore("useConfigStore", () => {
  const config = ref<Config>({});

  const read = async () => {
    let res = await get("/api/config");
    try {
      config.value = JSON.parse(res);
    } catch {
      config.value = {};
    }
    return config.value;
  };

  const write = useDebounceFn(async () => {
    await post("/api/config", config.value);
  }, 1000);

  read().then(() => {
    if (!config.value.Service?.length) {
      config.value.Service = [];
    }

    if (!config.value.FileService?.length) {
      config.value.FileService = [];
    }
  });

  return { config, read, write };
});

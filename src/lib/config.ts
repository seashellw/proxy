import { useDebounceFn } from "@vueuse/core";
import { defineStore } from "pinia";
import { ref } from "vue";
import { get, post } from "./fetch";
import { message } from "./message";

export interface HTTPSConfig {
  CertFile?: string;
  KeyFile?: string;
}
export interface ServiceConfig {
  Target?: string;
  Path?: string;
}
export interface RedirectConfig {
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
  Password?: string;
  Service?: ServiceConfig[];
  Redirect?: RedirectConfig[];
  FileService?: FileServiceConfig[];
  DynamicService?: DynamicServiceConfig;
  HTTPS?: HTTPSConfig;
}

const readPassword = () => {
  const password = localStorage.getItem("password");
  return password || "";
};

const writePassword = (password: string) => {
  localStorage.setItem("password", password);
};

const formatFileService = (config: Config) => {
  if (!config.FileService) {
    config.FileService = [];
  }
  config.FileService = config.FileService.map((item) => ({
    ...item,
    Dir: item.Dir?.trim() || "",
    Path: item.Path?.trim() || "",
  })).filter((item) => item.Dir);
};

const formatRedirect = (config: Config) => {
  if (!config.Redirect) {
    config.Redirect = [];
  }
  config.Redirect = config.Redirect.map((item) => ({
    ...item,
    Target: item.Target?.trim() || "",
    Path: item.Path?.trim() || "",
  })).filter((item) => item.Target);
};

const formatService = (config: Config) => {
  if (!config.Service) {
    config.Service = [];
  }
  config.Service = config.Service.map((item) => ({
    ...item,
    Target: item.Target?.trim() || "",
    Path: item.Path?.trim() || "",
  })).filter((item) => item.Target);
};

export const useConfigStore = defineStore("useConfigStore", () => {
  const config = ref<Config>({});

  const format = () => {
    formatService(config.value);
    formatRedirect(config.value);
    formatFileService(config.value);
  };

  const read = async () => {
    let res = await get("/api/config");
    config.value = res || {};
    format();
    return config.value;
  };

  const write = useDebounceFn(async () => {
    let password = readPassword();
    if (!password) {
      password = prompt("请输入密码") || "";
      writePassword(password);
    }
    config.value.Password = password;
    format();
    try {
      await post("/api/configSet", config.value);
    } catch (e) {
      console.error(e);
      message.error("写入配置失败");
      writePassword("");
      return;
    }
    message.success("写入配置成功");
  }, 500);

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

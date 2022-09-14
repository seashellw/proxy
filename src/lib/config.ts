import { useDebounceFn } from "@vueuse/core";
import { reactive } from "vue";
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
  Service?: ServiceConfig[];
  Redirect?: RedirectConfig[];
  FileService?: FileServiceConfig[];
  DynamicService?: DynamicServiceConfig;
  HTTPS?: HTTPSConfig;
}

const getPassword = () => {
  let password = localStorage.getItem("password") || "";
  if (!password) {
    password = (prompt("请输入密码") || "").trim();
    localStorage.setItem("password", password);
  }
  return password;
};

const pathSet = new Set<string>();

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

const config = reactive<Config>({});

const format = () => {
  formatService(config);
  formatRedirect(config);
  formatFileService(config);
};

const checkRepetition = () => {
  const { Service, Redirect, FileService } = config;
  pathSet.clear();
  let count = 0;
  const add = (item: { Path?: string }) => {
    pathSet.add(item.Path || "");
    count++;
  };
  Service?.forEach(add);
  Redirect?.forEach(add);
  FileService?.forEach(add);
  return pathSet.size === count;
};

const read = async () => {
  let password = getPassword();
  let res = await get(`/api/config?password=${encodeURIComponent(password)}`);
  Object.assign(config, res);
  format();
  return config;
};

const write = useDebounceFn(async () => {
  let password = getPassword();
  format();
  if (!checkRepetition()) {
    message.error("路径重复，禁止写入");
    return;
  }
  try {
    await post(
      `/api/configSet?password=${encodeURIComponent(password)}`,
      config
    );
  } catch (e) {
    console.error(e);
    message.error("写入配置失败");
    return;
  }
  message.success("写入配置成功");
}, 500);

export const init = async () => {
  await read();
  if (!config.Service?.length) {
    config.Service = [];
  }
  if (!config.FileService?.length) {
    config.FileService = [];
  }
};

export const useConfigStore = () => {
  return { config, read, write };
};

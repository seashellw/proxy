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

export interface StaticConfig {
  Path?: string;
  Dir?: string;
}

export interface ServerConfig {
  Service?: ServiceConfig[];
  Redirect?: RedirectConfig[];
  Static?: StaticConfig[];
  DynamicService?: DynamicServiceConfig;
  HTTPS?: HTTPSConfig;
}

export interface ConfigItem {
  id: string;
  value: string[];
}

interface Config {
  Service: ConfigItem[];
  Redirect: ConfigItem[];
  Static: ConfigItem[];
}

export const getId = (() => {
  let id = 0;
  return () => {
    id++;
    return id.toString();
  };
})();

const config = reactive<Config>({
  Service: [],
  Redirect: [],
  Static: [],
});

const getPassword = () => {
  let password = localStorage.getItem("password") || "";
  if (!password) {
    password = (prompt("请输入密码") || "").trim();
    localStorage.setItem("password", password);
  }
  return password;
};

let _config: any = null;

const read = async () => {
  let password = getPassword();
  let res: ServerConfig = await get(
    `/api/config?password=${encodeURIComponent(password)}`
  );
  if (!res) {
    return;
  }
  _config = res;
  config.Service =
    res.Service?.map((item: ServiceConfig) => ({
      id: getId(),
      value: [item.Path || "", item.Target || ""],
    })) || [];
  config.Redirect =
    res.Redirect?.map((item: RedirectConfig) => ({
      id: getId(),
      value: [item.Path || "", item.Target || ""],
    })) || [];
  config.Static =
    res.Static?.map((item: StaticConfig) => ({
      id: getId(),
      value: [item.Path || "", item.Dir || ""],
    })) || [];
};

const write = useDebounceFn(async () => {
  let password = getPassword();
  try {
    await post(`/api/configSet?password=${encodeURIComponent(password)}`, {
      ..._config,
      Service: config.Service?.map((item) => ({
        Path: item.value[0],
        Target: item.value[1],
      })),
      Redirect: config.Redirect?.map((item) => ({
        Path: item.value[0],
        Target: item.value[1],
      })),
      Static: config.Static?.map((item) => ({
        Path: item.value[0],
        Dir: item.value[1],
      })),
    });
  } catch (e) {
    console.error(e);
    message.error("写入配置失败");
    return;
  }
  message.success("写入配置成功");
}, 500);

export const useConfigList = () => {
  return {
    config,
    read,
    write,
  };
};

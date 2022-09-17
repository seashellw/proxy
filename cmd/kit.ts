import { which } from "zx";

export const isWin = () => {
  return process.platform.includes("win");
};

export let exeName = isWin() ? "proxy-server.exe" : "proxy-server";

export const existCmd = async (cmd: string) => {
  try {
    await which(cmd);
    return true;
  } catch {
    return false;
  }
};

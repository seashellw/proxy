import { $, which } from "zx";

export const isWin = () => {
  return process.platform.includes("win");
};

export let exeName = isWin() ? "proxy-server.exe" : "proxy-server";

export const existCmd = async (cmd) => {
  try {
    await which(cmd);
    return true;
  } catch {
    return false;
  }
};

export const install = async () => {
  if (!(await existCmd("pnpm"))) {
    console.log("install pnpm");
    if (isWin()) {
      await $`iwr https://get.pnpm.io/install.ps1 -useb | iex`;
    } else {
      await $`curl -fsSL https://get.pnpm.io/install.sh | sh -`;
    }
  }
  if (!(await existCmd("pm2"))) {
    console.log("install pm2");
    await $`pnpm add --global pm2`;
  }

  if (!((await existCmd("go")) || isWin())) {
    await $`curl -LO https://get.golang.org/$(uname)/go_installer && chmod +x go_installer && ./go_installer && rm go_installer`;
  }
};

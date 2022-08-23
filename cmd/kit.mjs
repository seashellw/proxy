export const isWin = () => {
  return process.platform.includes("win");
};

export let exeName = isWin() ? "proxy-server.exe" : "proxy-server";

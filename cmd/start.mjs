import { $ } from "zx";
import { isWin } from "./kit.mjs";

if (isWin()) {
  $`pm2 start ./proxy-server.exe`;
} else {
  $`pm2 start ./proxy-server`;
}

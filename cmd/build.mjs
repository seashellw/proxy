import { $ } from "zx";
import { isWin } from "./kit.mjs";

if (isWin()) {
  $`go build -o proxy-server.exe`;
} else {
  $`go build -o proxy-server`;
}

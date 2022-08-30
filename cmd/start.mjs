import { $ } from "zx";
import { exeName, install } from "./kit.mjs";

await install();

await $`pnpm pm2 start ./${exeName}`;

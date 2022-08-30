import { $ } from "zx";
import { exeName } from "./kit.mjs";

await $`pnpm pm2 start ./${exeName}`;

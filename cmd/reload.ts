import { $ } from "zx";
import { exeName } from "./kit";

await $`pnpm pm2 reload ./${exeName}`;

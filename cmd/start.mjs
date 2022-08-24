import { $ } from "zx";
import { exeName, install } from "./kit.mjs";

await install();

await $`pm2 start ./${exeName}`;

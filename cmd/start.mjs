import { $ } from "zx";
import { exeName } from "./kit.mjs";

await $`pm2 start ./${exeName}`;

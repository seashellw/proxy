import { $ } from "zx";
import { exeName, install } from "./kit.mjs";

await install();

await $`pnpm vue-tsc --noEmit && pnpm vite build`;
await $`go build -o ${exeName}`;

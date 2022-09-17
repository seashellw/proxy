import { $ } from "zx";
import { exeName } from "./kit.js";

await $`pnpm vue-tsc --noEmit && pnpm vite build`;
await $`go build -o ${exeName}`;

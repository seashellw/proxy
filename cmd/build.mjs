import { $ } from "zx";
import { exeName } from "./kit.mjs";

await $`pnpm vue-tsc --noEmit && pnpm vite build`;
await $`go build -o ${exeName}`;

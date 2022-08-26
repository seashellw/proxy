import { createDiscreteApi, darkTheme } from "naive-ui";

export let { message, dialog, notification, loadingBar } =
  createDiscreteApi(
    ["message", "dialog", "notification", "loadingBar"],
    {
      configProviderProps: {
        theme: darkTheme,
      },
    }
  );

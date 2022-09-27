<script setup lang="ts">
import { darkTheme, NConfigProvider, NLayout, NLayoutContent } from "naive-ui";
import { defineAsyncComponent } from "vue";
import { useConfigList } from "./lib/configList";

const { read } = useConfigList();
read();

const StaticList = defineAsyncComponent(
  () => import("./components/StaticList.vue")
);

const LayoutHeader = defineAsyncComponent(
  () => import("./components/LayoutHeader.vue")
);

const ServiceList = defineAsyncComponent(
  () => import("./components/ServiceList.vue")
);

const RedirectList = defineAsyncComponent(
  () => import("./components/RedirectList.vue")
);
</script>

<template>
  <NConfigProvider :theme="darkTheme">
    <NLayout style="height: 100vh">
      <LayoutHeader style="height: 4rem" />
      <NLayoutContent style="height: calc(100vh - 4rem)">
        <ServiceList />
        <RedirectList />
        <StaticList />
      </NLayoutContent>
    </NLayout>
  </NConfigProvider>
</template>

<style scoped></style>

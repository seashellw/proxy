<script setup lang="ts">
import {
  darkTheme,
  NConfigProvider,
  NLayout,
  NLayoutContent,
  zhCN,
} from "naive-ui";
import FileServiceList from "./components/FileServiceList.vue";
import LayoutHeader from "./components/LayoutHeader.vue";
import ServiceList from "./components/ServiceList.vue";
import RedirectList from "./components/RedirectList.vue";
</script>

<template>
  <NConfigProvider :theme="darkTheme" :locale="zhCN">
    <NLayout style="height: 100vh">
      <LayoutHeader style="height: 4rem" />
      <NLayoutContent style="height: calc(100vh - 4rem)">
        <ServiceList />
        <RedirectList />
        <FileServiceList />
      </NLayoutContent>
    </NLayout>
  </NConfigProvider>
</template>

<style scoped></style>

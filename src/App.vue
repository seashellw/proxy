<script setup lang="ts">
import { darkTheme, NConfigProvider, zhCN } from "naive-ui";
import FileServiceList from "./components/FileServiceList.vue";
</script>

<template>
  <NConfigProvider :theme="darkTheme" :locale="zhCN" class="app">
    <FileServiceList />
  </NConfigProvider>
</template>

<style scoped>
.app {
  width: 100vw;
  height: 100vh;
  overflow-x: hidden;
  overflow-y: auto;
}
</style>

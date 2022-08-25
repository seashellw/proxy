<script setup lang="ts">
import { NInput, NForm } from "naive-ui";
import { onMounted, ref, watch, watchEffect } from "vue";
import { Config, readConfig, writeConfig } from "./lib/config";
import FileService from "./components/FileService.vue";

const config = ref<Config>({});

onMounted(async () => {
  config.value = await readConfig();
  if (!config.value.Service?.length) {
    config.value.Service = [];
  }
  if (!config.value.FileService?.length) {
    config.value.FileService = [];
  }
});
</script>

<template>
  <div>
    <FileService
      :config="item"
      v-for="item in config.FileService"
      :key="item.Path"
    />
  </div>
</template>

<style scoped></style>

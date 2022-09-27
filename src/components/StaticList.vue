<script setup lang="ts">
import { toRefs } from "vue";
import { useConfigList } from "../lib/configList";
import InputTable from "./InputTable.vue";

const { config } = useConfigList();

let { Static } = toRefs(config);

const head = ["源地址", "本机路径"];
</script>

<template>
  <h3>静态资源服务</h3>
  <InputTable v-model:list="Static" :head="head" />
</template>

<style scoped></style>

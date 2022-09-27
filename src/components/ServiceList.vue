<script setup lang="ts">
import { toRefs } from "vue";
import { useConfigList } from "../lib/configList";
import InputTable from "./InputTable.vue";

const { config } = useConfigList();

let { Service } = toRefs(config);

const head = ["源地址", "目标地址"];
</script>

<template>
  <h3>代理服务</h3>
  <InputTable v-model:list="Service" :head="head" />
</template>

<style scoped></style>

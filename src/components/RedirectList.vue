<script setup lang="ts">
import { toRefs } from "vue";
import { useConfigList } from "../lib/configList";
import InputTable from "./InputTable.vue";

const { config } = useConfigList();

let { Redirect } = toRefs(config);

const head = ["源地址", "目标地址"];
</script>

<template>
  <h3>重定向服务</h3>
  <InputTable v-model:list="Redirect" :head="head" />
</template>

<style scoped></style>

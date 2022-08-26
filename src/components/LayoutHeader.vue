<script setup lang="ts">
import { Bolt } from "@vicons/tabler";
import { toRefs } from "@vueuse/core";
import { NButton, NIcon, NLayoutHeader } from "naive-ui";
import { useConfigStore } from "../lib/config";

const { write } = toRefs(useConfigStore());
const handleSave = () => {
  write.value();
};
</script>

<template>
  <NLayoutHeader class="header">
    <h1>服务器配置面板</h1>
    <NButton @click="handleSave" type="primary">
      <template #icon>
        <NIcon>
          <Bolt />
        </NIcon>
      </template>
      <span> 保存配置 </span>
    </NButton>
  </NLayoutHeader>
</template>

<style scoped>
h1 {
  margin: 0;
  flex-grow: 1;
  font-size: 1.1rem;
}
.header {
  display: flex;
  padding: 1rem;
  align-items: center;
}
</style>

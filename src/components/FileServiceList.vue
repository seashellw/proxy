<script setup lang="ts">
import { Minus, Plus } from "@vicons/tabler";
import { NButton, NIcon, NList, NListItem } from "naive-ui";
import { toRefs } from "vue";
import { useConfigStore } from "../lib/config";
import FileService from "./FileService.vue";

const { config } = toRefs(useConfigStore());

const handleAdd = () => {
  if (!config.value.FileService) {
    config.value.FileService = [];
  }
  config.value.FileService = [
    {
      Path: "",
      Dir: "",
    },
    ...config.value.FileService,
  ];
};

const handleRemove = (index: number) => {
  config.value.FileService?.splice(index, 1);
};
</script>

<template>
  <NList class="form-list" bordered>
    <template #header>
      <div class="form-header">
        <h2>静态文件服务</h2>
        <NButton @click="handleAdd">
          <template #icon>
            <NIcon>
              <Plus />
            </NIcon>
          </template>
        </NButton>
      </div>
    </template>
    <NListItem v-for="(_, index) in config.FileService" :key="index">
      <div class="form-list-item">
        <FileService :index="index" class="item-input" />
        <NButton @click="handleRemove(index)">
          <template #icon>
            <NIcon>
              <Minus />
            </NIcon>
          </template>
        </NButton>
      </div>
    </NListItem>
  </NList>
</template>

<style scoped>
h2 {
  margin: 0;
  flex-grow: 1;
  font-size: 1.2rem;
}
.form-list {
  margin: 1rem;
}

.form-header {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.form-list-item {
  display: flex;
  align-items: center;
}

.form-list-item .item-input {
  flex-grow: 1;
}
</style>

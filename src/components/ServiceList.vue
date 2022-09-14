<script setup lang="ts">
import { Minus, Plus } from "@vicons/tabler";
import { NButton, NIcon, NList, NListItem } from "naive-ui";
import { useConfigStore } from "../lib/config";
import Service from "./Service.vue";

const { config } = useConfigStore();

const handleAdd = () => {
  if (!config.Service) {
    config.Service = [];
  }
  config.Service.push({
    Path: "",
    Target: "",
  });
};

const handleRemove = (index: number) => {
  config.Service?.splice(index, 1);
};
</script>

<template>
  <NList class="form-list" bordered>
    <template #header>
      <div class="form-header">
        <h2>代理服务</h2>
        <NButton @click="handleAdd" type="info" circle>
          <template #icon>
            <NIcon>
              <Plus />
            </NIcon>
          </template>
        </NButton>
      </div>
    </template>
    <NListItem v-for="(_, index) in config.Service" :key="index">
      <div class="form-list-item">
        <Service :index="index" class="item-input" />
        <NButton @click="handleRemove(index)" type="warning" circle>
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
  gap: 1rem;
}

.form-list-item .item-input {
  flex-grow: 1;
}
</style>

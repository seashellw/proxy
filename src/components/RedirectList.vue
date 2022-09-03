<script setup lang="ts">
import { Minus, Plus } from "@vicons/tabler";
import { NButton, NIcon, NList, NListItem } from "naive-ui";
import { toRefs } from "vue";
import { useConfigStore } from "../lib/config";
import Redirect from "./Redirect.vue";

const { config } = toRefs(useConfigStore());

const handleAdd = () => {
  if (!config.value.Redirect) {
    config.value.Redirect = [];
  }
  config.value.Redirect.push({
    Path: "",
    Target: "",
  });
};

const handleRemove = (index: number) => {
  config.value.Service?.splice(index, 1);
};
</script>

<template>
  <NList class="form-list" bordered>
    <template #header>
      <div class="form-header">
        <h2>重定向服务</h2>
        <NButton @click="handleAdd" type="info">
          <template #icon>
            <NIcon>
              <Plus />
            </NIcon>
          </template>
        </NButton>
      </div>
    </template>
    <NListItem v-for="(_, index) in config.Redirect" :key="index">
      <div class="form-list-item">
        <Redirect :index="index" class="item-input" />
        <NButton @click="handleRemove(index)" type="warning">
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

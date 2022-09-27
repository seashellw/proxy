<script setup lang="ts">
import { Minus, Plus } from "@vicons/tabler";
import { useVModels } from "@vueuse/core";
import { NButton, NInput, NTable } from "naive-ui";
import { ConfigItem, getId } from "../lib/configList";

const props = defineProps<{
  head: string[];
  list: ConfigItem[];
}>();

const emit = defineEmits<{
  (key: string, value: any): void;
}>();

const { list } = useVModels(props, emit);

const handleRemove = (index: number) => {
  list.value?.splice(index, 1);
};

const handleCreate = () => {
  list.value?.push({
    id: getId(),
    value: ["", ""],
  });
};
</script>

<template>
  <NTable size="small">
    <thead>
      <tr>
        <th class="header-item" v-for="(item, index) in head" :key="index">
          {{ item }}
        </th>
        <th class="header-item action-col">
          <NButton quaternary size="small" circle @click="handleCreate">
            <template #icon>
              <Plus />
            </template>
          </NButton>
        </th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="(item, index) in list" :key="item.id">
        <td>
          <NInput v-model:value="item.value[0]" />
        </td>
        <td>
          <NInput v-model:value="item.value[1]" />
        </td>
        <td class="action-col">
          <NButton
            quaternary
            circle
            size="small"
            @click="() => handleRemove(index)"
          >
            <template #icon>
              <Minus />
            </template>
          </NButton>
        </td>
      </tr>
      <tr v-if="!list.length">
        <td class="header-item">无</td>
        <td class="header-item">无</td>
      </tr>
    </tbody>
  </NTable>
</template>

<style scoped>
.header-item {
  padding-left: 1rem;
}
.action-col {
  width: 2.5rem;
  text-align: center;
  padding: 0;
}
</style>

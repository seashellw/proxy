<script setup lang="ts">
import { useVModels } from "@vueuse/core";
import { FormItemRule, NForm, NInput } from "naive-ui";
import { ref } from "vue";
import { FileServiceConfig } from "../lib/config";
import { testPath } from "../lib/formRules";

const props = defineProps<{
  config: FileServiceConfig;
}>();

const emit = defineEmits<{
  (e: "update:config", value: FileServiceConfig): void;
}>();

const { config } = useVModels(props, emit);

const rules = ref({
  Path: {
    validator(_: FormItemRule, value: string) {
      if (!value) {
        return true;
      }
      if (!testPath(value)) {
        return new Error("路径不能以斜杠结尾");
      }
      return true;
    },
  },
  Dir: {
    validator(_: FormItemRule, value: string) {
      if (!value) {
        return new Error("必须填写目录路径");
      }
      return true;
    },
  },
});
</script>

<template>
  <NForm ref="formRef" :model="config" :rules="rules">
    <NInput v-model:value="config.Path" placeholder="源路径前缀" />
    <NInput v-model:value="config.Dir" placeholder="文件目录" />
  </NForm>
</template>

<style scoped></style>

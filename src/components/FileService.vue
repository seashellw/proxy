<script setup lang="ts">
import { ChevronsLeft } from "@vicons/tabler";
import { FormRules, NForm, NFormItem, NIcon, NInput } from "naive-ui";
import { ref } from "vue";
import { useConfigStore } from "../lib/config";
import { testPath } from "../lib/formRules";

defineProps<{
  index: number;
}>();

const {
  config: { FileService },
} = useConfigStore();

const rules = ref<FormRules>({
  Path: [
    {
      trigger: ["input", "blur"],
      validator(_, value: string) {
        if (!value) {
          return true;
        }
        if (!testPath(value)) {
          return new Error("路径不能以斜杠结尾");
        }
        return true;
      },
    },
  ],
  Dir: [
    {
      required: true,
      trigger: ["input", "blur"],
      validator(_, value: string) {
        if (!value) {
          return new Error("必须填写目录路径");
        }
        return true;
      },
    },
  ],
});
</script>

<template>
  <NForm
    ref="formRef"
    :model="FileService?.[index]"
    v-if="FileService?.[index]"
    :rules="rules"
    inline
    class="form"
  >
    <NFormItem path="Path" label="源路径前缀" class="form-item">
      <NInput class="input" v-model:value="FileService[index].Path" />
    </NFormItem>
    <NIcon class="icon" size="1.5rem">
      <ChevronsLeft />
    </NIcon>
    <NFormItem path="Dir" label="文件目录" class="form-item">
      <NInput class="input" v-model:value="FileService[index].Dir" />
    </NFormItem>
  </NForm>
</template>

<style scoped>
.form-item {
  flex-grow: 1;
}

.icon {
  align-self: center;
  margin-right: 1rem;
}
</style>

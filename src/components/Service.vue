<script setup lang="ts">
import { ChevronsRight } from "@vicons/tabler";
import { FormRules, NForm, NFormItem, NIcon, NInput } from "naive-ui";
import { ref, toRefs } from "vue";
import { useConfigStore } from "../lib/config";
import { testPath, testUrl } from "../lib/formRules";

defineProps<{
  index: number;
}>();

const { config } = toRefs(useConfigStore());

const rules = ref<FormRules>({
  Path: [
    {
      required: true,
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
  Target: [
    {
      required: true,
      trigger: ["input", "blur"],
      validator(_, value: string) {
        if (!value) {
          return new Error("必须填写目标路径");
        }
        if (!testUrl(value)) {
          return new Error("目标路径格式不正确");
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
    :model="config.Service?.[index]"
    v-if="config.Service?.[index]"
    :rules="rules"
    inline
    class="form"
  >
    <NFormItem path="Path" label="源路径前缀" class="form-item">
      <NInput class="input" v-model:value="config.Service[index].Path" />
    </NFormItem>
    <NIcon class="icon" size="1.5rem">
      <ChevronsRight />
    </NIcon>
    <NFormItem path="Target" label="目标路径前缀" class="form-item">
      <NInput class="input" v-model:value="config.Service[index].Target" />
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

<script setup lang="ts">
  import { Label } from "reka-ui";

  defineProps<{
    modelValue: string;
    type?: string;
    placeholder?: string;
    id?: string;
    label: string;
  }>();

  defineEmits<{
    "update:modelValue": [value: string];
  }>();
</script>

<template>
  <div class="flex flex-col gap-1.5">
    <Label :for="id" class="text-sm font-medium text-stone-600">
      {{ label }}
    </Label>
    <input
      :id="id"
      :type="type ?? 'text'"
      :value="modelValue"
      :placeholder="placeholder"
      class="input-field"
      v-bind="$attrs"
      @input="
        $emit('update:modelValue', ($event.target as HTMLInputElement).value)
      "
    />
  </div>
</template>

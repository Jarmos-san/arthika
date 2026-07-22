<script setup lang="ts">
  /**
   * @description Labeled text input with v-model support. All other attributes (e.g. `id`,
   * `type`, `placeholder`, `class`) fall through to the native `<input>`
   * element.
   */
  import { Label } from "reka-ui";
  import { useAttrs } from "vue";

  import type Props from "~/types/components/Input";

  defineOptions({ inheritAttrs: false });

  const props = defineProps<Props>();
  const model = defineModel<string>();
  const attrs = useAttrs();
</script>

<template>
  <div class="flex flex-col gap-1.5">
    <Label :for="attrs.id as string" class="text-sm font-medium text-stone-600">
      {{ props.label }}
    </Label>
    <input
      v-model="model"
      class="rounded-lg border border-stone-300 bg-stone-50 px-3.5 py-2.5 text-sm text-stone-800 placeholder:text-stone-400"
      v-bind="attrs"
    />
  </div>
</template>

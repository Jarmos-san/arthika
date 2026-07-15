<script setup lang="ts">
  import { Icon } from "@iconify/vue";
  import { CheckboxIndicator, CheckboxRoot } from "reka-ui";
  import { ref } from "vue";

  const checkboxOne = ref(true);
</script>

<template>
  <div class="flex flex-col gap-2.5">
    <label
      class="flex flex-row items-center gap-4 [&>.checkbox]:hover:bg-neutral-100"
    >
      <CheckboxRoot
        v-model="checkboxOne"
        class="flex h-5 w-5 appearance-none items-center justify-center rounded-md border bg-white shadow-sm outline-none focus-within:shadow-[0_0_0_2px_black] hover:bg-stone-50"
      >
        <CheckboxIndicator
          class="flex h-full w-full items-center justify-center rounded bg-white"
        >
          <Icon icon="radix-icons:check" class="text-grass9 h-5 w-5" />
        </CheckboxIndicator>
      </CheckboxRoot>
      <span class="text-sm text-stone-700 select-none dark:text-white"
        >Accept terms and conditions.</span
      >
    </label>
  </div>
</template>

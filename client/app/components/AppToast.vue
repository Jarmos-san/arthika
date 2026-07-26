<script setup lang="ts">
  import { ToastDescription, ToastRoot } from "reka-ui";

  interface Props {
    /** @description The message to display in the toast. */
    message: string;
  }

  interface Emits {
    /** @description Fired when the toast is dismissed (closed by timeout, swipe, or button). */
    dismiss: [];
  }

  const props = defineProps<Props>();
  const emit = defineEmits<Emits>();

  const onOpenChange = (open: boolean): void => {
    if (!open) {
      emit("dismiss");
    }
  };
</script>

<template>
  <ToastRoot class="app-toast" type="foreground" @update:open="onOpenChange">
    <ToastDescription class="text-sm font-bold text-stone-700">
      {{ props.message }}
    </ToastDescription>
  </ToastRoot>
</template>

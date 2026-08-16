<script setup lang="ts">
  import type Props from "~/types/components/AssetClassConfirmDialog";

  interface Emits {
    /** @description Fired when the user confirms the deletion. */
    confirmed: [];

    /** @description Fired when the dialog's open state changes. */
    "update:open": [open: boolean];
  }

  const props = defineProps<Props>();
  const emit = defineEmits<Emits>();

  const onOpenChange = (open: boolean): void => {
    emit("update:open", open);
  };

  const onConfirm = (): void => {
    emit("confirmed");
  };
</script>

<template>
  <AlertDialogRoot :open="props.open" @update:open="onOpenChange">
    <AlertDialogPortal>
      <AlertDialogOverlay class="fixed inset-0 z-50 bg-stone-900/40" />
      <AlertDialogContent
        class="fixed top-1/2 left-1/2 z-50 w-full max-w-md -translate-x-1/2 -translate-y-1/2 rounded-xl bg-white p-6 shadow-xl"
      >
        <AlertDialogTitle
          class="text-lg font-semibold tracking-tight text-stone-800"
        >
          Delete {{ props.asset?.name }}?
        </AlertDialogTitle>
        <AlertDialogDescription class="mt-1.5 text-sm text-stone-500">
          This removes the asset class from your tracking. You can't undo this.
        </AlertDialogDescription>

        <div class="mt-6 flex justify-end gap-3">
          <AlertDialogCancel as-child>
            <button class="btn-secondary" type="button">Cancel</button>
          </AlertDialogCancel>
          <AlertDialogAction as-child>
            <button class="btn-danger" type="button" @click="onConfirm">
              Delete
            </button>
          </AlertDialogAction>
        </div>
      </AlertDialogContent>
    </AlertDialogPortal>
  </AlertDialogRoot>
</template>

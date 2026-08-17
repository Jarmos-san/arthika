<script setup lang="ts">
  import type { AssetClass } from "~/openapi";

  interface Props {
    /** @description The asset class pending deletion, or `undefined` when closed. */
    asset: AssetClass;
  }

  interface Emits {
    /** @description Fired when the user confirms the deletion. */
    confirmed: [asset: AssetClass];

    /** @description Fired when the dialog's open state changes. */
    "update:open": [open: boolean];
  }

  const props = defineProps<Props>();
  const emit = defineEmits<Emits>();

  /**
   * @description Forwards reka-ui's open-state changes to the parent so it can unmount the
   * dialog when the user dismisses it.
   *
   * @param {boolean} open Whether the dialog is still open.
   */
  const onOpenChange = (open: boolean): void => {
    emit("update:open", open);
  };

  /**
   * @description Emits the pending asset class when the user clicks Delete, so the parent
   * can run the delete request.
   */
  const onConfirm = (): void => {
    emit("confirmed", props.asset);
  };
</script>

<template>
  <AlertDialogRoot :open="true" @update:open="onOpenChange">
    <AlertDialogPortal>
      <AlertDialogOverlay class="fixed inset-0 z-50 bg-stone-900/40" />
      <AlertDialogContent
        class="fixed top-1/2 left-1/2 z-50 w-full max-w-md -translate-x-1/2 -translate-y-1/2 rounded-xl bg-white p-6 shadow-xl"
      >
        <AlertDialogTitle
          class="text-lg font-semibold tracking-tight text-stone-800"
        >
          Delete {{ props.asset.name }}?
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

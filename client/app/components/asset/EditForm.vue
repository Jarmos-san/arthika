<script setup lang="ts">
  import { DialogOverlay, DialogPortal, DialogRoot } from "reka-ui";
  import { ref } from "vue";

  import { validateAssetClass, useToast } from "#imports";
  import type { AssetClass } from "~/openapi";
  import type { AssetClassErrors } from "~/types/utils/validators";

  interface Props {
    /** @description The asset class being edited. */
    asset: AssetClass;

    /** @description Whether the dialog is open. */
    open: boolean;
  }

  interface Emits {
    /** @description Fired when the dialog's open state changes. */
    "update:open": [open: boolean];
  }

  const props = defineProps<Props>();
  const emit = defineEmits<Emits>();
  const toast = useToast();

  /**
   * @description Forwards reka-ui's open-state changes to the parent.
   *
   * @param {boolean} open Whether the dialog is still open.
   */
  const onOpenChange = (open: boolean): void => {
    emit("update:open", open);
  };

  const errors = ref<AssetClassErrors>({ name: undefined });

  const isSubmitting = ref(false);

  /**
   * @description Submit button handler: submits the updated asset information and closes the
   * form modal.
   */
  const onSubmit = (): void => {
    // Validate the form input values
    errors.value = validateAssetClass(props.asset.name);
    if (errors.value.name) {
      return;
    }

    isSubmitting.value = true;

    try {
      const msg = `Edited ${props.asset.name}`;
      emit("update:open", false);
      toast.publish(msg);
    } finally {
      isSubmitting.value = false;
    }
  };

  /** @description Cancel button handler: closes the dialog. */
  const onClose = (): void => {
    emit("update:open", false);
  };
</script>

<template>
  <DialogRoot :open="props.open" @update:open="onOpenChange">
    <DialogPortal>
      <!-- The transparent overlay -->
      <DialogOverlay class="fixed inset-0 z-50 bg-stone-900/40" />

      <!-- The actual content of the modal -->
      <DialogContent
        class="fixed top-1/2 left-1/2 z-50 w-full max-w-md -translate-x-1/2 -translate-y-1/2 rounded-xl bg-white p-6 shadow-xl"
      >
        <!-- The modal title -->
        <DialogTitle
          class="text-lg font-semibold tracking-tight text-stone-800"
        >
          Edit {{ props.asset.name }}
        </DialogTitle>

        <!-- The modal description -->
        <DialogDescription class="mt-1.5 text-sm text-stone-500">
          Update the details for {{ props.asset.name }} below. Changes are saved
          when you submit the form.
        </DialogDescription>

        <form @submit.prevent="onSubmit" class="mt-5 flex flex-col gap-5">
          <div>
            <!-- Asset name -->
            <Input
              id="asset-name"
              aria-describedby="asset-name-error"
              label="Name"
              placeholder="e.g. Equities"
            />
            <p
              v-if="errors.name"
              id="asset-name-errors"
              class="mt-1.5 text-xs text-red-600"
            >
              {{ errors.name }}
            </p>
          </div>

          <Textarea
            id="asset-description"
            label="Description"
            placeholder="What belongs in this asset class?"
            :rows="3"
          />
        </form>

        <!-- The button groups -->
        <div class="mt-6 flex justify-end gap-3">
          <!-- The cancel button -->
          <button class="btn-secondary" type="button" @click="onClose">
            Cancel
          </button>

          <!-- The submit button -->
          <button
            class="btn-primary"
            type="submit"
            @submit.prevent="onSubmit"
            @click="onSubmit"
          >
            Submit
          </button>
        </div>
      </DialogContent>
    </DialogPortal>
  </DialogRoot>
</template>

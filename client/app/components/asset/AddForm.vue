<script setup lang="ts">
  import { ref, reactive } from "vue";

  import { useToast, validateAssetClass } from "#imports";
  import type { AssetClass, AssetClassInput } from "~/openapi";
  import type { AssetClassErrors } from "~/types/utils/validators";

  interface Emits {
    /** @description Fired with the created asset class after a successful save. */
    saved: [asset: AssetClass];
  }

  const emit = defineEmits<Emits>();
  const toast = useToast();

  /** @description Whether the "Add asset class" button is open. */
  const isOpen = ref(false);

  /** @description Form fields for the new asset class. */
  const asset = reactive<AssetClassInput>({
    description: undefined,
    name: "",
  });

  /** @description Per-field validation errors; `undefined` means the field is valid. */
  const errors = ref<AssetClassErrors>({ name: undefined });

  /** @description True while the create request is in flight; disables the Save button. */
  const isSubmitting = ref(false);

  /** @description Clears the form fields and validation errors. */
  const resetForm = (): void => {
    isOpen.value = false;
    asset.name = "";
    asset.description = "";
    errors.value = {
      name: undefined,
    };
  };

  /** @description Cancel button handler: closes the dialog and resets the form. */
  const onClose = (): void => {
    resetForm();
  };

  /**
   * @description Validates the form, creates the asset class via POST /api/assets, emits
   * `saved` with the created asset, and shows a toast on success. No-op when
   * validation fails.
   */
  const onSubmit = async (): Promise<void> => {
    // Validate the form input values.
    errors.value = validateAssetClass(asset.name);
    if (errors.value.name) {
      return;
    }

    // Disable the save button when the request is in flight.
    isSubmitting.value = true;

    // Attempt to pass the client-side input values to the server for persistence,
    // create a toast notification and then reset the form before resetting the
    // submission state
    try {
      const created = await $fetch<AssetClass>("/api/assets", {
        body: { description: asset.description || undefined, name: asset.name },
        method: "POST",
      });
      emit("saved", created);
      toast.publish(`Created new asset class: "${asset.name}"`);
      resetForm();
    } finally {
      isSubmitting.value = false;
    }
  };
</script>

<template>
  <!-- Button to add a new asset class -->
  <button class="btn-positive" type="button" @click="isOpen = true">
    Add asset class
  </button>

  <!-- Form modal -->
  <DialogRoot v-model:open="isOpen">
    <DialogPortal>
      <DialogOverlay class="fixed inset-0 z-50 bg-stone-900/40" />
      <DialogContent
        class="fixed top-1/2 left-1/2 z-50 w-full max-w-md -translate-x-1/2 -translate-y-1/2 rounded-xl bg-white p-6 shadow-xl"
      >
        <!-- Form title -->
        <DialogTitle
          class="text-lg font-semibold tracking-tight text-stone-800"
        >
          Add asset class
        </DialogTitle>

        <form class="mt-5 flex flex-col gap-5" @submit.prevent="onSubmit">
          <div>
            <!-- Asset name -->
            <Input
              id="asset-name"
              v-model="asset.name"
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

          <!-- Asset description -->
          <Textarea
            id="asset-description"
            v-model="asset.description"
            label="Description"
            placeholder="What belongs in this asset class?"
            :rows="3"
          />
        </form>

        <div class="mt-5 flex justify-end gap-3">
          <!-- Cancel button -->
          <button class="btn-secondary" type="button" @click="onClose">
            Cancel
          </button>

          <!-- Save button -->
          <button
            class="btn-positive"
            type="submit"
            :disabled="isSubmitting"
            @click="onSubmit"
          >
            {{ isSubmitting ? "Saving..." : "Save" }}
          </button>
        </div>
      </DialogContent>
    </DialogPortal>
  </DialogRoot>
</template>

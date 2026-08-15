<script setup lang="ts">
  import { ref, reactive } from "vue";

  import { useToast } from "#imports";
  import type { AssetClass, AssetClassInput } from "~/openapi";
  import type { AssetClassErrors } from "~/types/utils/validators";
  import { validateAssetClass } from "~/utils/validators";

  interface Emits {
    saved: [asset: AssetClass];
  }

  // Initial state of the add asset class form modal set to "not open"
  const isOpen = ref(false);

  // Reactive object containing the asset details
  const asset = reactive<AssetClassInput>({
    description: "",
    name: "",
  });

  // Reactive object of the error values to be shown on the form if validation fails
  const errors = ref<AssetClassErrors>({ name: undefined });

  const isSubmitting = ref(false);

  const emit = defineEmits<Emits>();
  const toast = useToast();

  const resetForm = (): void => {
    isOpen.value = false;
    asset.name = "";
    asset.description = "";
    errors.value = {
      name: undefined,
    };
  };

  const onClose = (): void => {
    resetForm();
  };

  const onSubmit = async (): Promise<void> => {
    errors.value = validateAssetClass(asset.name);
    if (errors.value.name) {
      return;
    }

    isSubmitting.value = true;

    try {
      const created = await $fetch<AssetClass>("/api/assets", {
        body: { description: asset.description || undefined, name: asset.name },
        method: "POST",
      });
      emit("saved", created);
      toast.publish(`Created new asset class: ${asset.name}`);
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

        <form action="submit" class="mt-5 flex flex-col gap-5">
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

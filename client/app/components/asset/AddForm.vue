<script setup lang="ts">
  import { ref, reactive } from "vue";

  import type { AssetClassInput } from "~/openapi";
  import type { AssetClassErrors } from "~/types/utils/validators";

  // Initial state of the add asset class form modal set to "not open"
  const isOpen = ref(false);

  // Reactive object containing the asset details
  const asset = reactive<AssetClassInput>({
    description: "",
    name: "",
  });

  // Reactive object of the error values to be shown on the form if validation fails
  const errors = ref<AssetClassErrors>({ name: undefined });
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

        <!-- Cancel add asset action -->
        <div class="mt-5 flex justify-end gap-3">
          <DialogClose as-child>
            <button class="btn-secondary" type="button">Cancel</button>
          </DialogClose>
        </div>
      </DialogContent>
    </DialogPortal>
  </DialogRoot>
</template>

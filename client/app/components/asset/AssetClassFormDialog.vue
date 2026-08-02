<script setup lang="ts">
  import { computed, ref, watch } from "vue";

  import useAssetClasses from "~/composables/useAssetClasses";
  import type Props from "~/types/components/AssetClassFormDialog";
  import type {
    AssetClass,
    AssetClassInput,
  } from "~/types/composables/useAssetClasses";
  import type { AssetClassErrors } from "~/types/utils/validators";
  import { validateAssetClass } from "~/utils/validators";

  interface Emits {
    /** @description Fired with the saved asset class after a successful create or update. */
    saved: [asset: AssetClass];

    /** @description Fired when the dialog's open state changes. */
    "update:open": [open: boolean];
  }

  const props = defineProps<Props>();
  const emit = defineEmits<Emits>();

  const { createAssetClass, updateAssetClass } = useAssetClasses();

  const name = ref("");
  const description = ref("");
  const errors = ref<AssetClassErrors>({ name: undefined });
  const isSubmitting = ref(false);

  const dialogTitle = computed<string>(() => {
    if (props.asset === undefined) {
      return "Add asset class";
    }

    return "Edit asset class";
  });

  // Reset the form with the asset being edited (or blank for a new asset class)
  watch(
    () => props.open,
    (isOpen) => {
      if (!isOpen) {
        return;
      }
      errors.value = { name: undefined };
      if (props.asset === undefined) {
        name.value = "";
        description.value = "";
        return;
      }
      name.value = props.asset.name;
      description.value = props.asset.description;
    },
  );

  const onOpenChange = (open: boolean): void => {
    emit("update:open", open);
  };

  const buildInput = (): AssetClassInput => ({
    description: description.value.trim(),
    name: name.value.trim(),
  });

  const persist = async (
    input: Readonly<AssetClassInput>,
  ): Promise<AssetClass> => {
    if (props.asset !== undefined) {
      return await updateAssetClass(props.asset.id, input);
    }

    return await createAssetClass(input);
  };

  const onSubmit = async (): Promise<void> => {
    errors.value = validateAssetClass(name.value);
    if (errors.value.name !== undefined) {
      return;
    }

    isSubmitting.value = true;
    try {
      const saved = await persist(buildInput());
      emit("saved", saved);
      emit("update:open", false);
    } finally {
      isSubmitting.value = false;
    }
  };
</script>

<template>
  <DialogRoot :open="props.open" @update:open="onOpenChange">
    <DialogPortal>
      <DialogOverlay class="fixed inset-0 z-50 bg-stone-900/40" />
      <DialogContent
        class="fixed top-1/2 left-1/2 z-50 w-full max-w-md -translate-x-1/2 -translate-y-1/2 rounded-xl bg-white p-6 shadow-xl"
      >
        <DialogTitle
          class="text-lg font-semibold tracking-tight text-stone-800"
        >
          {{ dialogTitle }}
        </DialogTitle>

        <form class="mt-5 flex flex-col gap-5" @submit.prevent="onSubmit">
          <div>
            <Input
              id="asset-name"
              v-model="name"
              aria-describedby="asset-name-error"
              label="Name"
              placeholder="e.g. Equities"
            />
            <p
              v-if="errors.name"
              id="asset-name-error"
              class="mt-1.5 text-xs text-red-600"
            >
              {{ errors.name }}
            </p>
          </div>

          <Textarea
            id="asset-description"
            v-model="description"
            label="Description"
            placeholder="What belongs in this asset class?"
            :rows="3"
          />

          <div class="mt-1 flex justify-end gap-3">
            <DialogClose as-child>
              <button
                class="btn-secondary"
                type="button"
                :disabled="isSubmitting"
              >
                Cancel
              </button>
            </DialogClose>
            <button class="btn-positive" type="submit" :disabled="isSubmitting">
              {{ isSubmitting ? "Saving..." : "Save changes" }}
            </button>
          </div>
        </form>
      </DialogContent>
    </DialogPortal>
  </DialogRoot>
</template>

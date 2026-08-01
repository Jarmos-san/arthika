<script setup lang="ts">
  import { useHead } from "nuxt/app";
  import { ref } from "vue";

  import AssetClassConfirmDialog from "~/components/asset/AssetClassConfirmDialog.vue";
  import AssetClassFormDialog from "~/components/asset/AssetClassFormDialog.vue";
  import AssetClassTable from "~/components/asset/AssetClassTable.vue";
  import useAssetClasses from "~/composables/useAssetClasses";
  import useToast from "~/composables/useToast";
  import type { AssetClass } from "~/types/composables/useAssetClasses";

  useHead({ title: "Dashboard" });

  const { assetClasses, isMutating, removeAssetClass } = useAssetClasses();
  const { publish } = useToast();

  const isFormOpen = ref(false);
  const isConfirmOpen = ref(false);
  const editingAsset = ref<AssetClass | undefined>(undefined);
  const pendingDelete = ref<AssetClass | undefined>(undefined);

  /** @description Opens the form dialog in create mode. */
  const openCreate = (): void => {
    editingAsset.value = undefined;
    isFormOpen.value = true;
  };

  /**
   * @description Opens the form dialog pre-filled with the given asset class.
   *
   * @param {AssetClass} asset - The asset class to edit.
   */
  const openEdit = (asset: AssetClass): void => {
    editingAsset.value = asset;
    isFormOpen.value = true;
  };

  /**
   * @description Opens the confirm dialog for the given asset class.
   *
   * @param {AssetClass} asset - The asset class pending deletion.
   */
  const requestDelete = (asset: AssetClass): void => {
    pendingDelete.value = asset;
    isConfirmOpen.value = true;
  };

  /**
   * @description Announces the result of a create or update via a toast.
   *
   * @param {AssetClass} saved - The asset class that was created or updated.
   */
  const onSaved = (saved: AssetClass): void => {
    if (editingAsset.value !== undefined) {
      publish(`Updated ${saved.name}`);
      return;
    }
    publish(`Added ${saved.name}`);
  };

  /** @description Removes the pending asset class and announces it via a toast. */
  const onConfirmedDelete = async (): Promise<void> => {
    if (pendingDelete.value === undefined) {
      return;
    }
    const { name } = pendingDelete.value;
    await removeAssetClass(pendingDelete.value.id);
    isConfirmOpen.value = false;
    pendingDelete.value = undefined;
    publish(`Deleted ${name}`);
  };
</script>

<template>
  <div class="mx-auto w-full max-w-4xl px-6 py-10">
    <header class="flex items-end justify-between gap-4">
      <div>
        <h1 class="text-3xl font-semibold tracking-tight text-stone-800">
          Asset Classes
        </h1>
        <p class="mt-1 text-sm text-stone-500">
          {{ assetClasses.length }}
          {{ assetClasses.length === 1 ? "asset class" : "asset classes" }}
          tracked
        </p>
      </div>
      <button
        class="btn-positive"
        type="button"
        :disabled="isMutating"
        @click="openCreate"
      >
        Add asset class
      </button>
    </header>

    <main class="mt-8">
      <AssetClassTable
        :asset-classes="assetClasses"
        @add="openCreate"
        @edit="openEdit"
        @delete="requestDelete"
      />
    </main>

    <AssetClassFormDialog
      v-model:open="isFormOpen"
      :asset="editingAsset"
      @saved="onSaved"
    />

    <AssetClassConfirmDialog
      v-model:open="isConfirmOpen"
      :asset="pendingDelete"
      @confirmed="onConfirmedDelete"
    />
  </div>
</template>

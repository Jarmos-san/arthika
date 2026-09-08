<script setup lang="ts">
  import DeleteIcon from "@iconify-vue/material-symbols/delete";
  import EditOutlineRoundedIcon from "@iconify-vue/material-symbols/edit-outline-rounded";
  import { ref } from "vue";

  import useDeleteAsset from "~/composables/useDeleteAsset";
  import type { AssetClass } from "~/openapi";

  interface Props {
    /** @description Asset classes to render as ledger rows. */
    assetClasses: AssetClass[] | undefined;
  }

  interface Emits {
    /** @description Fired with the deleted ID after a successful delete operation. */
    deleted: [id: string];
  }

  defineProps<Props>();
  const emit = defineEmits<Emits>();

  const {
    pendingDelete,
    requestDelete,
    onDialogOpenChange,
    onDeleteConfirmed,
  } = useDeleteAsset({ onDeleted: (id) => emit("deleted", id) });

  /** @description The list of header elements to render on top of the table. */
  const headers = ["CLASS", "DESCRIPTION", "ACTIONS"];

  const pendingEdit = ref<AssetClass | undefined>(undefined);

  /**
   * @description Clears the pending edit when the dialog reports it closed.
   *
   * @param {boolean} open Boolean flag to denote if the edit form is open or
   *   not.
   */
  const onEditDialogOpenChange = (open: boolean): void => {
    if (!open) {
      pendingEdit.value = undefined;
    }
  };
</script>

<template>
  <div class="rounded-xl border border-stone-200/60 bg-white shadow-lg">
    <div class="overflow-x-auto">
      <table
        v-if="assetClasses && assetClasses?.length > 0"
        class="w-full text-sm"
      >
        <!-- Table header -->
        <thead>
          <tr class="border-b border-stone-200 text-left">
            <th
              v-for="(header, index) in headers"
              :key="index"
              class="px-6 py-3.5 font-mono text-xs font-medium tracking-wider text-stone-400"
              :class="{ 'text-right': index === header.length - 1 }"
            >
              {{ header }}
            </th>
          </tr>
        </thead>

        <!-- Table body -->
        <tbody class="divide-y divide-stone-200">
          <tr
            v-for="assetClass in assetClasses"
            :key="assetClass.id"
            class="transition-colors duration-150 hover:bg-stone-50/60"
          >
            <!-- Name of the class -->
            <td class="px-6 py-4 font-medium text-stone-800">
              {{ assetClass.name }}
            </td>

            <!-- Description of the class -->
            <td class="px-6 py-4 text-stone-500">
              {{ assetClass.description }}
            </td>

            <!-- Action buttons -->
            <td class="p-4">
              <div class="flex justify-start gap-1">
                <!-- Edit button -->
                <button
                  class="btn-ghost"
                  type="button"
                  :aria-label="`Edit asset class ${assetClass.name}`"
                  @click="pendingEdit = assetClass"
                >
                  <EditOutlineRoundedIcon height="1rem" />
                </button>

                <!-- Delete button -->
                <button
                  class="btn-ghost text-red-600 hover:text-red-700"
                  type="button"
                  :aria-label="`Delete asset class ${assetClass.name}`"
                  @click="requestDelete(assetClass)"
                >
                  <DeleteIcon height="1rem" />
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>

      <div
        v-else
        class="m-6 flex flex-col items-center rounded-xl border-2 border-dashed border-stone-200 p-10 text-center"
      >
        <h2 class="text-base font-semibold tracking-tight text-stone-800">
          No asset classes tracked yet
        </h2>
        <p class="mt-1.5 text-sm text-stone-500">
          Add your first one to start tracking.
        </p>
        <button class="btn-positive mt-5" type="button">
          Add asset class
        </button>
      </div>
    </div>

    <!-- Delete confirmation dialog -->
    <AssetDeleteConfirmation
      v-if="pendingDelete"
      :asset="pendingDelete"
      @confirmed="onDeleteConfirmed"
      @update:open="onDialogOpenChange"
    />

    <!-- Edit asset form modal -->
    <AssetEditForm
      v-if="pendingEdit"
      :asset="pendingEdit"
      :open="true"
      @update:open="onEditDialogOpenChange"
    />
  </div>
</template>

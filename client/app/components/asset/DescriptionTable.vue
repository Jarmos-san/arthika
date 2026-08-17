<script setup lang="ts">
  import DeleteIcon from "@iconify-vue/material-symbols/delete";
  import EditOutlineRoundedIcon from "@iconify-vue/material-symbols/edit-outline-rounded";
  import { ref } from "vue";

  import { useToast } from "#imports";
  import type { AssetClass, DeleteAsset204 } from "~/openapi";

  import DeleteConfirmation from "./DeleteConfirmation.vue";

  interface Props {
    /** @description Asset classes to render as ledger rows. */
    assetClasses: AssetClass[] | undefined;
  }

  interface Emits {
    /** @description Fired with the deleted ID after a successful delete operation. */
    deleted: [id: string];
  }

  const props = defineProps<Props>();
  const emit = defineEmits<Emits>();
  const toast = useToast();

  /** @description The list of header elements to render on top of the table. */
  const headers = ref(["CLASS", "DESCRIPTION", "ACTIONS"]);

  /** @description ID of the row currently being deleted; `undefined` when idle. */
  const pendingId = ref<string | undefined>(undefined);

  /**
   * @description Asset class awaiting delete confirmation; `undefined` when no dialog is
   * open.
   */
  const pendingDelete = ref<AssetClass>();

  /**
   * @description Deletes the asset class via DELETE /api/assets/{id}, shows a toast
   * notification on success or failure, and emits `deleted` with the ID so the
   * parent can refresh the list. Guards against double-submits via
   * `pendingId`.
   *
   * @param {AssetClass} assetClass Represents an asset class (e.g., "Equity",
   *   "Debt Bonds", etc).
   */
  const deleteAssetHandler = async (assetClass: AssetClass): Promise<void> => {
    pendingId.value = assetClass.id;

    try {
      await $fetch<DeleteAsset204>(`/api/assets/${assetClass.id}`, {
        method: "DELETE",
      });
      toast.publish(`Deleted asset class: "${assetClass.name}"`);
      emit("deleted", assetClass.id);
    } catch {
      toast.publish("Failed to delete asset class. Please try again.");
    } finally {
      pendingId.value = undefined;
    }
  };

  /**
   * @description Open the confirmation dialog for the given asset class.
   *
   * @param {AssetClassl} assetClass An instance of the asset class
   *   representation.
   */
  const requestDelete = (assetClass: AssetClass): void => {
    pendingDelete.value = assetClass;
  };

  /**
   * @description Sync dialog visibility; clears the pending asset when closed.
   *
   * @param {boolean} open Boolean value representing the visibility state of
   *   the modal.
   */
  const onDialogOpenChange = (open: boolean): void => {
    if (!open) {
      pendingDelete.value = undefined;
    }
  };

  /**
   * @description Runs the delete after the user confirms; closes and resets state.
   *
   * @param {AssetClass} assetClass The asset class the user confirmed deleting,
   *   received from the confirmation dialog's `confirmed` emit.
   */
  const onDeleteConfirmed = async (assetClass: AssetClass): Promise<void> => {
    if (!assetClass) {
      return;
    }

    pendingDelete.value = undefined;
    await deleteAssetHandler(assetClass);
  };
</script>

<template>
  <div class="rounded-xl border border-stone-200/60 bg-white shadow-lg">
    <div class="overflow-x-auto">
      <table
        v-if="props.assetClasses ? props.assetClasses.length > 0 : undefined"
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
            v-for="assetClass in props.assetClasses"
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
                <button class="btn-ghost" type="button">
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
    <DeleteConfirmation
      v-if="pendingDelete"
      :asset="pendingDelete"
      @confirmed="onDeleteConfirmed"
      @update:open="onDialogOpenChange"
    />
  </div>
</template>

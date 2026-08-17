import { ref, type Ref } from "vue";

import { useToast } from "#imports";
import type { AssetClass } from "~/openapi";

interface UseDeleteAssetOptions {
  /** @description Called with the deleted ID after a successful delete. */
  onDeleted: (id: string) => void;
}

interface UseDeleteAsset {
  /**
   * @description Asset class awaiting delete confirmation; `undefined` when no dialog is
   * open.
   */
  pendingDelete: Ref<AssetClass | undefined>;

  /** @description Opens the confirmation dialog for the given asset class. */
  requestDelete: (assetClass: Readonly<AssetClass>) => void;

  /** @description Syncs dialog visibility; clears the pending asset when closed. */
  onDialogOpenChange: (open: boolean) => void;

  /** @description Runs the delete after the user confirms; closes and resets state. */
  onDeleteConfirmed: (assetClass: Readonly<AssetClass>) => Promise<void>;
}

/**
 * @description Composable for the asset-class delete flow: confirmation dialog state,
 * double-submit guard, DELETE /api/assets/{id} request, and toast
 * notifications. Reports success through `onDeleted` so the caller can refresh
 * its data.
 *
 * @param {UseDeleteAssetOptions} options Callback invoked with the deleted ID.
 *
 * @returns {UseDeleteAsset} Dialog state and handlers for the delete flow.
 */
const useDeleteAsset = (
  options: Readonly<UseDeleteAssetOptions>,
): UseDeleteAsset => {
  const toast = useToast();

  /** @description ID of the row currently being deleted; `undefined` when idle. */
  const pendingId = ref<string | undefined>(undefined);

  /**
   * @description Asset class awaiting delete confirmation; `undefined` when no dialog is
   * open.
   */
  const pendingDelete = ref<AssetClass | undefined>(undefined);

  /**
   * @description Deletes the asset class via DELETE /api/assets/{id}, shows a toast on
   * success or failure, and invokes `onDeleted` with the ID. Guards against
   * double-submits via `pendingId`.
   *
   * @param {AssetClass} assetClass Represents an asset class (e.g., "Equity",
   *   "Debt Bonds", etc).
   */
  const deleteAssetHandler = async (
    assetClass: Readonly<AssetClass>,
  ): Promise<void> => {
    if (pendingId.value !== undefined) {
      return;
    }

    pendingId.value = assetClass.id;

    try {
      await $fetch(`/api/assets/${assetClass.id}`, { method: "DELETE" });
      toast.publish(`Deleted asset class: "${assetClass.name}"`);
      options.onDeleted(assetClass.id);
    } catch {
      toast.publish("Failed to delete asset class. Please try again.");
    } finally {
      pendingId.value = undefined;
    }
  };

  const requestDelete = (assetClass: Readonly<AssetClass>): void => {
    pendingDelete.value = assetClass;
  };

  const onDialogOpenChange = (open: boolean): void => {
    if (!open) {
      pendingDelete.value = undefined;
    }
  };

  const onDeleteConfirmed = async (
    assetClass: Readonly<AssetClass>,
  ): Promise<void> => {
    pendingDelete.value = undefined;
    await deleteAssetHandler(assetClass);
  };

  return {
    onDeleteConfirmed,
    onDialogOpenChange,
    pendingDelete,
    requestDelete,
  };
};

export default useDeleteAsset;

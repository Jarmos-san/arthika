import type { AssetClass } from "~/types/composables/useAssetClasses";

/**
 * @description Modal dialog for creating and editing asset classes. The dialog owns the form
 * state and validation, and emits the saved asset class to the caller.
 */
export default interface Props {
  /** @description The asset class being edited, or `undefined` when creating a new one. */
  asset: AssetClass | undefined;
  /** @description Whether the dialog is open. */
  open: boolean;
}

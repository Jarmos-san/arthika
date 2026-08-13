import type { AssetClass } from "~/openapi";

/**
 * @description Modal dialog confirming the deletion of an asset class. Emits `confirmed`
 * when the user chooses to delete.
 */
export default interface Props {
  /** @description The asset class pending deletion, or `undefined` when closed. */
  asset: AssetClass | undefined;

  /** @description Whether the dialog is open. */
  open: boolean;
}

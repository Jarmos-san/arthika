/** @description Types for asset classes managed by the `useAssetClasses` composable. */

/** @description A tracked asset class as shown on the dashboard. */
interface AssetClass {
  /** @description Unique identifier for the asset class. */
  id: string;

  /** @description Display name of the asset class (e.g. "Equities"). */
  name: string;

  /** @description What belongs in this asset class. */
  description: string;
}

/** @description Input payload for creating or updating an asset class. */
interface AssetClassInput {
  /** @description Display name of the asset class. */
  name: string;

  /** @description What belongs in this asset class. */
  description: string;
}

export type { AssetClass, AssetClassInput };

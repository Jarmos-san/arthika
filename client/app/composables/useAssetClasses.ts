import { useFetch } from "nuxt/app";

import type { AssetClass } from "~/openapi";

/**
 * @description Result type of the `useFetch` call backing this composable. Derived via
 * `ReturnType` so the exposed shape stays in sync with Nuxt's `AsyncData` type
 * without hand-written drift.
 */
type Result = ReturnType<typeof useFetch<AssetClass[]>>;

/**
 * @description Composable for fetching asset classes from the server-side `/api/assets`
 * endpoint. Wraps Nuxt's `useFetch` so the pag can consume reactive asset-class
 * data with SSR support: the request runs on the server during initial render,
 * and the session cookie is forwarded automatically for same-origin requests.
 *
 * @returns {Pick<Result, "data" | "error" | "refresh" | "status">} An object
 *   containing the details returned by the server-side application.
 */
const useAssetClasses = (): Pick<
  Result,
  "data" | "error" | "refresh" | "status"
> => {
  const { data, error, refresh, status } =
    useFetch<AssetClass[]>("/api/assets");

  return { data, error, refresh, status };
};

export default useAssetClasses;

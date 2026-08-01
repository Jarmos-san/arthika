import { ref, type Ref } from "vue";

import type {
  AssetClass,
  AssetClassInput,
} from "~/types/composables/asset-classes";

/** @description Simulated network latency for stub CRUD operations, in milliseconds. */
const STUB_LATENCY_MS = 400;

/** @description Sentinel returned by `findIndex` when no match is found. */
const INDEX_NOT_FOUND = -1;

/** @description Number of items removed per splice call. */
const REMOVE_COUNT = 1;

/** @description Seed data so the dashboard shows the design on first load. */
const SEED_CLASSES: readonly AssetClass[] = [
  {
    description: "Public stocks, ETFs, and index funds",
    id: "seed-equities",
    name: "Equities",
  },
  {
    description: "Bonds, treasuries, and cash equivalents",
    id: "seed-fixed-income",
    name: "Fixed income",
  },
  {
    description: "Bitcoin, Ethereum, and other digital assets",
    id: "seed-crypto",
    name: "Crypto",
  },
];

/** @description In-memory list of asset classes. Module scope so all callers share one list. */
const assetClasses = ref<AssetClass[]>(
  SEED_CLASSES.map((seed: Readonly<AssetClass>) => ({
    description: seed.description,
    id: seed.id,
    name: seed.name,
  })),
);

/** @description Whether a stub CRUD operation is currently in flight. */
const isMutating = ref(false);

/**
 * @description Waits for the simulated latency before running the mutation, keeping
 * `isMutating` set while the operation is in flight.
 *
 * @template Result - The return type of the mutation.
 * @param {() => Promise<Result> | Result} mutate - The mutation to run after
 *   the latency.
 *
 * @returns {Promise<Result>} The result of the mutation.
 */
const withStubLatency = async <Result>(
  mutate: () => Promise<Result> | Result,
): Promise<Result> => {
  isMutating.value = true;
  try {
    // oxlint-disable-next-line promise/avoid-new -- a timer promise simulates network latency
    await new Promise((resolve) => {
      setTimeout(resolve, STUB_LATENCY_MS);
    });
    return await mutate();
  } finally {
    isMutating.value = false;
  }
};

/**
 * @description Creates a new asset class and appends it to the shared list.
 *
 * @param {Readonly<AssetClassInput>} input - The name and description of the
 *   new asset class.
 *
 * @returns {Promise<AssetClass>} The created asset class.
 */
const createAssetClass = async (
  input: Readonly<AssetClassInput>,
): Promise<AssetClass> => {
  const created = await withStubLatency(() => {
    const assetClass: AssetClass = {
      description: input.description,
      id: crypto.randomUUID(),
      name: input.name,
    };
    assetClasses.value.push(assetClass);
    return assetClass;
  });
  return created;
};

/**
 * @description Updates the matching asset class in place.
 *
 * @param {string} id - The id of the asset class to update.
 * @param {Readonly<AssetClassInput>} input - The new name and description.
 *
 * @returns {Promise<AssetClass>} The updated asset class.
 */
const updateAssetClass = async (
  id: string,
  input: Readonly<AssetClassInput>,
): Promise<AssetClass> => {
  const updated = await withStubLatency(() => {
    const target = assetClasses.value.find(
      (item: Readonly<AssetClass>) => item.id === id,
    );
    if (target === undefined) {
      throw new Error(`Unknown asset class: ${id}`);
    }
    target.name = input.name;
    target.description = input.description;
    return {
      description: target.description,
      id: target.id,
      name: target.name,
    };
  });
  return updated;
};

/**
 * @description Removes the matching asset class from the shared list. Silently ignores
 * unknown ids.
 *
 * @param {string} id - The id of the asset class to remove.
 *
 * @returns {Promise<void>} Resolves once the operation completes.
 */
const removeAssetClass = async (id: string): Promise<void> => {
  await withStubLatency(() => {
    const index = assetClasses.value.findIndex(
      (item: Readonly<AssetClass>) => item.id === id,
    );
    if (index !== INDEX_NOT_FOUND) {
      assetClasses.value.splice(index, REMOVE_COUNT);
    }
  });
};

/** @description Return type for the {@link useAssetClasses} composable. */
interface UseAssetClasses {
  /** @description Reactive list of asset classes. */
  assetClasses: Ref<AssetClass[]>;
  /** @description Creates a new asset class and appends it to the list. */
  createAssetClass: (input: Readonly<AssetClassInput>) => Promise<AssetClass>;
  /** @description Whether a stub CRUD operation is currently in flight. */
  isMutating: Ref<boolean>;
  /** @description Removes the matching asset class from the list. */
  removeAssetClass: (id: string) => Promise<void>;
  /** @description Updates the matching asset class in place. */
  updateAssetClass: (
    id: string,
    input: Readonly<AssetClassInput>,
  ) => Promise<AssetClass>;
}

/**
 * @description Composable providing stub CRUD operations for asset classes. State lives at
 * module scope (like `useToast`) so every caller shares one list. The stub
 * latency keeps loading states visible; swapping the mutations for `$fetch`
 * calls later is a one-function change.
 *
 * @returns {UseAssetClasses} An object with `assetClasses`, `createAssetClass`,
 *   `updateAssetClass`, `removeAssetClass`, and `isMutating`.
 */
const useAssetClasses = (): UseAssetClasses => ({
  assetClasses,
  createAssetClass,
  isMutating,
  removeAssetClass,
  updateAssetClass,
});

export type { UseAssetClasses };
export default useAssetClasses;

import { beforeEach, describe, expect, it, vi } from "vitest";

import type { UseAssetClasses } from "~/composables/useAssetClasses";
import type { AssetClass } from "~/types/composables/useAssetClasses";

vi.setConfig({ testTimeout: 10_000 });

/** @description Number of seed classes shown on first load. */
const SEED_CLASS_COUNT = 3;

/** @description Number of classes after creating one. */
const CREATED_CLASS_COUNT = 4;

/** @description Number of classes after removing one. */
const REMAINING_CLASS_COUNT = 2;

/** @description Seed id used to target update operations in tests. */
const EQUITIES_SEED_ID = "seed-equities";

/** @description Seed id used to target removal operations in tests. */
const CRYPTO_SEED_ID = "seed-crypto";

describe("useAssetClasses", () => {
  let mod: { default: () => UseAssetClasses };

  beforeEach(async () => {
    vi.resetModules();
    mod = await import("~/composables/useAssetClasses");
  });

  it("starts with the seeded asset classes", () => {
    expect.hasAssertions();

    const { assetClasses } = mod.default();

    expect(assetClasses.value).toHaveLength(SEED_CLASS_COUNT);
  });

  it("createAssetClass appends and returns the new asset class", async () => {
    expect.hasAssertions();

    const { assetClasses, createAssetClass } = mod.default();
    const input = { description: "Green energy producers", name: "Renewables" };

    const created = await createAssetClass(input);

    expect(created).toMatchObject(input);
    expect(created.id).not.toBe("");
    expect(assetClasses.value).toHaveLength(CREATED_CLASS_COUNT);
    expect(
      assetClasses.value.some(
        (item: Readonly<AssetClass>) => item.id === created.id,
      ),
    ).toBeTruthy();
  });

  it("createAssetClass toggles isMutating while in flight", async () => {
    expect.hasAssertions();

    const { createAssetClass, isMutating } = mod.default();

    const pending = createAssetClass({
      description: "",
      name: "Renewables",
    });

    expect(isMutating.value).toBeTruthy();

    await pending;

    expect(isMutating.value).toBeFalsy();
  });

  it("updateAssetClass mutates the matching asset class in place", async () => {
    expect.hasAssertions();

    const { assetClasses, updateAssetClass } = mod.default();
    const input = { description: "Broad market exposure", name: "Stocks" };

    const updated = await updateAssetClass(EQUITIES_SEED_ID, input);

    expect(updated).toMatchObject(input);
    expect(updated.id).toBe(EQUITIES_SEED_ID);
    expect(assetClasses.value).toHaveLength(SEED_CLASS_COUNT);
    expect(
      assetClasses.value.some(
        (item: Readonly<AssetClass>) => item.id === EQUITIES_SEED_ID,
      ),
    ).toBeTruthy();
  });

  it("removeAssetClass deletes the matching asset class", async () => {
    expect.hasAssertions();

    const { assetClasses, removeAssetClass } = mod.default();

    await removeAssetClass(CRYPTO_SEED_ID);

    expect(assetClasses.value).toHaveLength(REMAINING_CLASS_COUNT);
    expect(
      assetClasses.value.some(
        (item: Readonly<AssetClass>) => item.id === CRYPTO_SEED_ID,
      ),
    ).toBeFalsy();
  });
});

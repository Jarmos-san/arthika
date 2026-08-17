import { registerEndpoint } from "@nuxt/test-utils/runtime";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { useToast } from "#imports";
import useDeleteAsset from "~/composables/useDeleteAsset";
import type { AssetClass } from "~/openapi";

vi.setConfig({ testTimeout: 10_000 });

describe("composables/useDeleteAsset", () => {
  const assetClass: AssetClass = {
    description: "Ownership shares in public or private companies.",
    id: "367f7bc7-9bfc-4e4d-a7c9-8bdec07f5c9f",
    name: "Equity",
  };

  let endpointCleanups: (() => void)[];

  beforeEach(() => {
    endpointCleanups = [];
    useToast().messages.value = [];
  });

  afterEach(() => {
    for (const cleanup of endpointCleanups) {
      cleanup();
    }

    endpointCleanups = [];
  });

  const rememberEndpoint = (
    url: string,
    // oxlint-disable-next-line typescript/prefer-readonly-parameter-types, eslint/no-magic-numbers
    options: Parameters<typeof registerEndpoint>[1],
  ): void => {
    endpointCleanups.push(registerEndpoint(url, options));
  };

  it("onDeleteConfirmed deletes the asset and reports success", async () => {
    expect.hasAssertions();

    rememberEndpoint(`/api/assets/${assetClass.id}`, {
      handler: () => ({}),
      method: "DELETE",
    });

    const onDeleted = vi.fn<() => void>();
    const { onDeleteConfirmed, pendingDelete, requestDelete } = useDeleteAsset({
      onDeleted,
    });

    requestDelete(assetClass);
    await onDeleteConfirmed(assetClass);

    expect(onDeleted).toHaveBeenCalledWith(assetClass.id);
    expect(useToast().messages.value).toContain(
      `Deleted asset class: "${assetClass.name}"`,
    );
    expect(pendingDelete.value).toBeUndefined();
  });

  it("onDeleteConfirmed publishes an error toast and does not report success on failure", async () => {
    expect.hasAssertions();

    rememberEndpoint(`/api/assets/${assetClass.id}`, {
      handler: (): Response =>
        Response.json({ message: "Internal server error" }, { status: 500 }),
      method: "DELETE",
    });

    const onDeleted = vi.fn<() => void>();
    const { onDeleteConfirmed } = useDeleteAsset({ onDeleted });

    await onDeleteConfirmed(assetClass);

    expect(onDeleted).not.toHaveBeenCalled();
    expect(useToast().messages.value).toContain(
      "Failed to delete asset class. Please try again.",
    );
  });

  it("guards against double-submits while a delete is in flight", async () => {
    expect.hasAssertions();

    const handler = vi.fn<() => object>(() => ({}));

    rememberEndpoint(`/api/assets/${assetClass.id}`, {
      handler,
      method: "DELETE",
    });

    const onDeleted = vi.fn<() => void>();
    const { onDeleteConfirmed } = useDeleteAsset({ onDeleted });

    const first = onDeleteConfirmed(assetClass);
    const second = onDeleteConfirmed(assetClass);

    await Promise.all([first, second]);

    expect(handler).toHaveBeenCalledOnce();
    expect(onDeleted).toHaveBeenCalledOnce();
  });
});

import { registerEndpoint } from "@nuxt/test-utils/runtime";
import { describe, expect, it, vi } from "vitest";

import { useAssetClasses } from "#imports";
import type { AssetClass } from "~/openapi";

vi.setConfig({ testTimeout: 10_000 });

describe("composables/useAssetClasses", () => {
  const sampleAssetClasses: AssetClass[] = [
    {
      description: "Ownership shares in public or private companies.",
      id: "367f7bc7-9bfc-4e4d-a7c9-8bdec07f5c9f",
      name: "Equity",
    },
    {
      description:
        "Bonds and debt securities issued by governments or corporations.",
      id: "cf05abe5-9901-4b72-a1eb-19981636edd1",
      name: "Fixed Income",
    },
  ];

  it("fetches asset classes on success", async () => {
    expect.hasAssertions();

    registerEndpoint("/api/assets", {
      handler: () => sampleAssetClasses,
      method: "GET",
    });

    const result = useAssetClasses();

    await vi.waitFor(() => {
      expect(result.status.value).toBe("success");
    });

    expect(result.data.value).toStrictEqual(sampleAssetClasses);
  });
});

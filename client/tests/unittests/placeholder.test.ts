import { describe, expect, it, vi } from "vitest";

vi.setConfig({ testTimeout: 10_000 });

describe("placeholder", () => {
  it("passes", () => {
    expect.hasAssertions();

    expect(true).toBeTruthy();
  });
});

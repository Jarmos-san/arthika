import { createPinia, setActivePinia } from "pinia";
import { describe, expect, it, vi } from "vitest";

import useAuthStore from "../../app/stores/auth";

vi.setConfig({ testTimeout: 10_000 });

describe("auth store", () => {
  it("is not authenticated by default", () => {
    expect.hasAssertions();
    setActivePinia(createPinia());

    const store = useAuthStore();
    expect(store.isAuthenticated).toBeFalsy();
  });

  it("sets isAuthenticated to true after login", () => {
    expect.hasAssertions();
    setActivePinia(createPinia());

    const store = useAuthStore();
    store.login("user@example.com", "password123");
    expect(store.isAuthenticated).toBeTruthy();
  });
});

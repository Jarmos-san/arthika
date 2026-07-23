import { createPinia, setActivePinia } from "pinia";
import { describe, expect, it, vi } from "vitest";

import useAuthStore from "../../app/stores/auth";

vi.setConfig({ testTimeout: 10_000 });

// Mock $fetch globally since it is injected by Nuxt at runtime.
const mockFetch = vi.fn<() => Promise<unknown>>();
vi.stubGlobal("$fetch", mockFetch);

describe("auth store", () => {
  it("is not authenticated by default", () => {
    expect.hasAssertions();
    setActivePinia(createPinia());

    const store = useAuthStore();
    expect(store.isAuthenticated).toBeFalsy();
  });

  it("sets isAuthenticated to true after successful login", async () => {
    expect.hasAssertions();
    setActivePinia(createPinia());

    mockFetch.mockResolvedValueOnce({
      email: "user@example.com",
      id: "550e8400-e29b-41d4-a716-446655440000",
    });

    const store = useAuthStore();
    await store.login("user@example.com", "password123");

    expect(store.isAuthenticated).toBeTruthy();
    expect(store.user).toStrictEqual({
      email: "user@example.com",
      id: "550e8400-e29b-41d4-a716-446655440000",
    });
    expect(mockFetch).toHaveBeenCalledWith("/api/users/login", {
      body: { email: "user@example.com", password: "password123" },
      method: "POST",
    });
  });

  it("throws on failed login", async () => {
    expect.hasAssertions();
    setActivePinia(createPinia());

    mockFetch.mockRejectedValueOnce(new Error("401 Unauthorized"));

    const store = useAuthStore();
    await expect(
      store.login("user@example.com", "wrongpassword"),
    ).rejects.toThrow("401 Unauthorized");
    expect(store.isAuthenticated).toBeFalsy();
  });

  it("fetchUser sets user on success", async () => {
    expect.hasAssertions();
    setActivePinia(createPinia());

    mockFetch.mockResolvedValueOnce({
      email: "user@example.com",
      id: "550e8400-e29b-41d4-a716-446655440000",
    });

    const store = useAuthStore();
    await store.fetchUser();

    expect(store.isAuthenticated).toBeTruthy();
    expect(store.user).toStrictEqual({
      email: "user@example.com",
      id: "550e8400-e29b-41d4-a716-446655440000",
    });
  });

  it("fetchUser clears state on failure", async () => {
    expect.hasAssertions();
    setActivePinia(createPinia());

    mockFetch.mockRejectedValueOnce(new Error("401 Unauthorized"));

    const store = useAuthStore();
    await store.fetchUser();

    expect(store.isAuthenticated).toBeFalsy();
    expect(store.user).toBeUndefined();
  });
});

import { describe, expect, it, vi } from "vitest";
import type { Ref } from "vue";
import { ref } from "vue";

vi.setConfig({ testTimeout: 10_000 });

declare global {
  // eslint-disable-next-line typescript/no-explicit-any
  var $fetch: ReturnType<typeof vi.fn>;
}

const stateMap = new Map<string, Ref<unknown>>();

const getOrCreateState = (key: string, init?: () => unknown): Ref<unknown> => {
  if (!stateMap.has(key)) {
    const defaultValue = init !== undefined ? init() : undefined;
    stateMap.set(key, ref<unknown>(defaultValue));
  }

  const state = stateMap.get(key);

  if (state === undefined) {
    return ref<unknown>(undefined);
  }

  return state;
};

vi.mock("#app", () => ({
  useAsyncData: (
    _key: string,
    handler: () => unknown,
  ): { data: Ref<unknown> } => {
    const promise = handler();
    if (promise instanceof Promise) {
      void promise;
    }
    const stateRef: Ref<unknown> =
      stateMap.get("auth-user") ?? ref<unknown>(undefined);
    return { data: stateRef };
  },
  useState: (key: string, init?: () => unknown): Ref<unknown> =>
    getOrCreateState(key, init),
}));

vi.stubGlobal("$fetch", vi.fn());

import useAuth, {
  STATUS_CONFLICT,
  STATUS_NETWORK_ERROR,
  STATUS_OK,
  STATUS_UNAUTHORIZED,
  STATUS_UNPROCESSABLE,
} from "../../../app/composables/useAuth";

describe("useAuth", () => {
  it("exports status constants", () => {
    expect.hasAssertions();
    expect(STATUS_OK).toBe(200);
    expect(STATUS_UNAUTHORIZED).toBe(401);
    expect(STATUS_CONFLICT).toBe(409);
    expect(STATUS_UNPROCESSABLE).toBe(422);
    expect(STATUS_NETWORK_ERROR).toBe(0);
  });

  it("is not authenticated by default", () => {
    expect.hasAssertions();
    $fetch.mockRejectedValue(new Error("Unauthorized"));
    stateMap.clear();

    const { isAuthenticated } = useAuth();

    expect(isAuthenticated.value).toBeFalsy();
  });

  it("sets isAuthenticated to true after login", async () => {
    expect.hasAssertions();
    $fetch.mockResolvedValue({
      email: "user@example.com",
      id: "550e8400-e29b-41d4-a716-446655440000",
    });
    stateMap.clear();

    const { isAuthenticated, login } = useAuth();

    await login("user@example.com", "password123");

    expect(isAuthenticated.value).toBeTruthy();
  });
});
